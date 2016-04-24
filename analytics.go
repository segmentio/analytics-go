package analytics

import (
	"io"
	"io/ioutil"

	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Version of the client.
const Version = "3.0.0"

type Client interface {
	io.Closer

	Enqueue(Message) error
}

type client struct {
	Config
	key string

	// This channel is where the `Enqueue` method writes messages so they can be
	// picked up and pushed by the backend goroutine taking care of applyinh the
	// batching rules.
	msgs chan interface{}

	// These two channels are used to synchronize the client shutting down when
	// `Close` is called.
	// The first channel is closed to signal the backend goroutine that it has
	// to stop, then the second one is closed by the backend goroutine to signal
	// that it has finished flushing all queued messages.
	quit     chan struct{}
	shutdown chan struct{}
}

func New(writeKey string) Client {
	c, _ := NewWithConfig(writeKey, Config{})
	return c
}

func NewWithConfig(writeKey string, config Config) (cli Client, err error) {
	if err = config.validate(); err != nil {
		return
	}

	c := &client{
		Config:   makeConfig(config),
		key:      writeKey,
		msgs:     make(chan interface{}, 100),
		quit:     make(chan struct{}),
		shutdown: make(chan struct{}),
	}

	go c.loop()

	cli = c
	return
}

func (c *client) Enqueue(msg Message) (err error) {
	if err = msg.validate(); err != nil {
		return
	}

	defer func() {
		// When the `msgs` channel is closed writing to it will trigger a panic.
		// To avoid letting the panic propagate to the caller we recover from it
		// and instead report that the client has been closed and shouldn't be
		// used anymore.
		if recover() != nil {
			err = io.EOF
		}
	}()

	c.msgs <- msg.serializable(c.UID(), c.Now())
	return
}

// Close and flush metrics.
func (c *client) Close() (err error) {
	defer func() {
		// Always recover, a panic could be raised if c.quit was closed which
		// means the Close method was called more than once.
		if recover() != nil {
			err = io.EOF
		}
	}()
	close(c.quit)
	<-c.shutdown
	return
}

// Send batch request.
func (c *client) send(msgs []interface{}) {
	if len(msgs) == 0 {
		return
	}

	b, err := json.Marshal(batch{
		MessageId: c.uid(),
		SentAt:    formatTime(c.now()),
		Messages:  msgs,
		Context:   defaultContext,
	})

	if err != nil {
		c.errorf("marshalling mesages - %s", err)
		return
	}

	for i := 0; i < 10; i++ {
		if err := c.upload(b); err == nil {
			break
		}
		time.Sleep(c.RetryAfter(i))
	}
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

	res, err := (&http.client{Transport: c.Transport}).Do(req)

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

	c.logf("response %s: %d – %s", res.Status, res.StatusCode, body)
}

// Batch loop.
func (c *client) loop() {
	defer close(c.shutdown)

	var msgs []interface{}
	var tick = time.NewTicker(c.Interval)
	defer tick.Stop()

	for {
		select {
		case msg := <-c.msgs:
			c.debugf("buffer (%d/%d) %v", len(msgs), c.BatchSize, msg)
			msgs = append(msgs, msg)

			if len(msgs) == c.BatchSize {
				c.debugf("exceeded %d messages – flushing", c.BatchSize)
				c.send(msgs)
				msgs = nil
			}

		case <-tick.C:
			if len(msgs) > 0 {
				c.debugf("interval reached - flushing %d", len(msgs))
				c.send(msgs)
				msgs = nil
			} else {
				c.debugf("interval reached – nothing to send")
			}

		case <-c.quit:
			c.debugf("exit requested – draining messages")

			// Drain the msg channel, we have to close it first so no more
			// messages can be pushed and otherwise the loop would never end.
			close(c.msgs)
			for msg := range c.msgs {
				c.debugf("buffer (%d/%d) %v", len(msgs), c.BatchSize, msg)
				msgs = append(msgs, msg)
			}

			c.debugf("exit requested – flushing %d", len(msgs))
			c.send(msgs)
			c.debugf("exit")
			return
		}
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
