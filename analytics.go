package analytics

import (
	"io"
	"io/ioutil"
	"sync"

	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Version of the client.
const Version = "3.0.0"

// This interface is the main API exposed by the analytics package.
// Values that satsify this interface are returned by the client constructors
// provided by the package and provide a way to send messages via the HTTP API.
type Client interface {
	io.Closer

	// Queues a message to be sent by the client when the conditions for a batch
	// upload are met.
	// This is the main method you'll be using, a typical flow would look like
	// this:
	//
	//	client := analytics.New(writeKey)
	//	...
	//	client.Enqueue(analytics.Track{ ... })
	//	...
	//	client.Close()
	//
	// The method returns an error if the message queue not be queued, which
	// happens if the client was already closed at the time the method was
	// called or if the message was malformed.
	Enqueue(Message) error
}

type client struct {
	Config
	key string

	// This channel is where the `Enqueue` method writes messages so they can be
	// picked up and pushed by the backend goroutine taking care of applying the
	// batching rules.
	msgs chan Message

	// These two channels are used to synchronize the client shutting down when
	// `Close` is called.
	// The first channel is closed to signal the backend goroutine that it has
	// to stop, then the second one is closed by the backend goroutine to signal
	// that it has finished flushing all queued messages.
	quit     chan struct{}
	shutdown chan struct{}

	// This HTTP client is used to send requests to the backend, it uses the
	// HTTP transport provided in the configuration.
	http http.Client
}

// Instantiate a new client that uses the write key passed as first argument to
// send messages to the backend.
// The client is created with the default configuration.
func New(writeKey string) Client {
	// Here we can ignore the error because the default config is always valid.
	c, _ := NewWithConfig(writeKey, Config{})
	return c
}

// Instantiate a new client that uses the write key and configuration passed as
// arguments to send messages to the backend.
// The function will return an error if the configuration contained impossible
// values (like a negative flush interval for example).
// When the function returns an error the returned client will always be nil.
func NewWithConfig(writeKey string, config Config) (cli Client, err error) {
	if err = config.validate(); err != nil {
		return
	}

	c := &client{
		Config:   makeConfig(config),
		key:      writeKey,
		msgs:     make(chan Message, 100),
		quit:     make(chan struct{}),
		shutdown: make(chan struct{}),
		http: http.Client{
			Transport: config.Transport,
		},
	}

	go c.loop()

	cli = c
	return
}

func (c *client) Enqueue(msg Message) (err error) {
	if err = msg.validate(); err != nil {
		return
	}

	var id = c.uid()
	var ts = c.now()

	switch m := msg.(type) {
	case Alias:
		m.Type = "alias"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m

	case Group:
		m.Type = "group"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m

	case Identify:
		m.Type = "identify"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m

	case Page:
		m.Type = "page"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m

	case Screen:
		m.Type = "screen"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m

	case Track:
		m.Type = "track"
		m.MessageId = makeMessageId(m.MessageId, id)
		m.Timestamp = makeTimestamp(m.Timestamp, ts)
		msg = m
	}

	defer func() {
		// When the `msgs` channel is closed writing to it will trigger a panic.
		// To avoid letting the panic propagate to the caller we recover from it
		// and instead report that the client has been closed and shouldn't be
		// used anymore.
		if recover() != nil {
			err = ErrClosed
		}
	}()

	c.msgs <- msg
	return
}

// Close and flush metrics.
func (c *client) Close() (err error) {
	defer func() {
		// Always recover, a panic could be raised if `c`.quit was closed which
		// means the method was called more than once.
		if recover() != nil {
			err = ErrClosed
		}
	}()
	close(c.quit)
	<-c.shutdown
	return
}

// Asynchronously send a batched requests.
func (c *client) sendAsync(msgs []message, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.send(msgs)
	}()
}

// Send batch request.
func (c *client) send(msgs []message) {
	const attempts = 10

	if len(msgs) == 0 {
		return
	}

	b, err := json.Marshal(batch{
		MessageId: c.uid(),
		SentAt:    c.now(),
		Messages:  msgs,
		Context:   c.DefaultContext,
	})

	if err != nil {
		c.errorf("marshalling mesages - %s", err)
		c.notifyFailure(msgs, err)
		return
	}

	for i := 0; i != attempts; i++ {
		if err = c.upload(b); err == nil {
			c.notifySucess(msgs)
			return
		}

		// Wait for either a retry timeout or the client to be closed.
		select {
		case <-time.After(c.RetryAfter(i)):
		case <-c.quit:
			c.errorf("%d messages dropped because they failed to be sent and the client was closed", len(msgs))
			return
		}
	}

	c.errorf("%d messages dropped because they failed to be sent after %d attempts", len(msgs), attempts)
	c.notifyFailure(msgs, err)
}

// Upload serialized batch message.
func (c *client) upload(b []byte) error {
	url := c.Endpoint + "/v1/batch"
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		c.errorf("creating request - %s", err)
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(b)))
	req.SetBasicAuth(c.key, "")

	res, err := c.http.Do(req)

	if err != nil {
		c.errorf("sending request - %s", err)
		return err
	}

	defer res.Body.Close()
	c.report(res)

	return nil
}

// Report on response body.
func (c *client) report(res *http.Response) {
	if res.StatusCode < 400 {
		c.debugf("response %s", res.Status)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.errorf("reading response body - %s", err)
		return
	}

	c.logf("response %s: %d – %s", res.Status, res.StatusCode, string(body))
}

// Batch loop.
func (c *client) loop() {
	defer close(c.shutdown)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	tick := time.NewTicker(c.Interval)
	defer tick.Stop()

	mq := messageQueue{
		maxBatchSize:  c.BatchSize,
		maxBatchBytes: c.maxBatchBytes(),
	}

	for {
		select {
		case msg := <-c.msgs:
			c.push(&mq, msg, wg)

		case <-tick.C:
			c.flush(&mq, wg)

		case <-c.quit:
			c.debugf("exit requested – draining messages")

			// Drain the msg channel, we have to close it first so no more
			// messages can be pushed and otherwise the loop would never end.
			close(c.msgs)
			for msg := range c.msgs {
				c.push(&mq, msg, wg)
			}

			c.flush(&mq, wg)
			c.debugf("exit")
			return
		}
	}
}

func (c *client) push(q *messageQueue, m Message, wg *sync.WaitGroup) {
	var msg message
	var err error

	if msg, err = makeMessage(m); err != nil {
		if c.Callback != nil {
			c.Callback.Failure(m, err)
		}
		return
	}

	c.debugf("buffer (%d/%d) %v", len(q.pending), c.BatchSize, m)

	if msgs := q.push(msg); msgs != nil {
		c.debugf("exceeded messages batch limit with batch of %d messages – flushing", len(msgs))
		c.sendAsync(msgs, wg)
	}
}

func (c *client) flush(q *messageQueue, wg *sync.WaitGroup) {
	if msgs := q.flush(); msgs != nil {
		c.debugf("flushing %d messages", len(msgs))
		c.sendAsync(msgs, wg)
	}
}

func (c *client) debugf(format string, args ...interface{}) {
	if c.Verbose {
		c.logf(format, args...)
	}
}

func (c *client) logf(format string, args ...interface{}) {
	c.Logger.Logf(format, args...)
}

func (c *client) errorf(format string, args ...interface{}) {
	c.Logger.Errorf(format, args...)
}

func (c *client) notifySucess(msgs []message) {
	if c.Callback != nil {
		for _, m := range msgs {
			c.Callback.Success(m.msg)
		}
	}
}

func (c *client) notifyFailure(msgs []message, err error) {
	if c.Callback != nil {
		for _, m := range msgs {
			c.Callback.Failure(m.msg, err)
		}
	}
}

func (c *client) maxBatchBytes() int {
	b, _ := json.Marshal(batch{
		MessageId: c.uid(),
		SentAt:    c.now(),
		Context:   c.DefaultContext,
	})
	return maxBatchBytes - len(b)
}
