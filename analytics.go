package analytics

import (
	"io"
	"io/ioutil"
	"sync"

	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/segmentio/backo-go"
	"github.com/xtgo/uuid"
)

// Version of the client.
const Version = "3.0.0"

// Endpoint for the Segment API.
const Endpoint = "https://api.segment.io"

// DefaultContext of message batches.
var DefaultContext = map[string]interface{}{
	"library": map[string]interface{}{
		"name":    "analytics-go",
		"version": Version,
	},
}

// Backoff policy.
var Backo = backo.DefaultBacko()

// Client which batches messages and flushes at the given Interval or
// when the Size limit is exceeded. Set Verbose to true to enable
// logging output.
type Client struct {
	Endpoint string
	// Interval represents the duration at which messages are flushed. It may be
	// configured only before any messages are enqueued.
	Interval  time.Duration
	Size      int
	Logger    Logger
	Verbose   bool
	Transport http.RoundTripper
	key       string
	msgs      chan interface{}
	quit      chan struct{}
	shutdown  chan struct{}
	uid       func() string
	now       func() time.Time
	once      sync.Once
}

// New client with write key.
func New(key string) *Client {
	c := &Client{
		Endpoint:  Endpoint,
		Interval:  5 * time.Second,
		Size:      250,
		Logger:    newDefaultLogger(),
		Verbose:   false,
		Transport: http.DefaultTransport,
		key:       key,
		msgs:      make(chan interface{}, 100),
		quit:      make(chan struct{}),
		shutdown:  make(chan struct{}),
		now:       time.Now,
		uid:       uid,
	}

	return c
}

func (c *Client) Enqueue(msg Message) (err error) {
	if err = msg.validate(); err != nil {
		return
	}
	c.once.Do(c.startLoop)
	c.msgs <- msg.serializable(c.uid(), c.now())
	return
}

func (c *Client) startLoop() {
	go c.loop()
}

// Close and flush metrics.
func (c *Client) Close() (err error) {
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
func (c *Client) send(msgs []interface{}) {
	if len(msgs) == 0 {
		return
	}

	b, err := json.Marshal((batch{
		Messages: msgs,
		Context:  DefaultContext,
	}).serializable(c.uid(), c.now()))

	if err != nil {
		c.errorf("marshalling mesages - %s", err)
		return
	}

	for i := 0; i < 10; i++ {
		if err := c.upload(b); err == nil {
			break
		}
		Backo.Sleep(i)
	}
}

// Upload serialized batch message.
func (c *Client) upload(b []byte) error {
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

	res, err := (&http.Client{Transport: c.Transport}).Do(req)

	if err != nil {
		c.errorf("sending request - %s", err)
		return err
	}

	defer res.Body.Close()
	c.report(res)

	return nil
}

// Report on response body.
func (c *Client) report(res *http.Response) {
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
func (c *Client) loop() {
	defer close(c.shutdown)

	var msgs []interface{}
	var tick = time.NewTicker(c.Interval)
	defer tick.Stop()

	for {
		select {
		case msg := <-c.msgs:
			c.debugf("buffer (%d/%d) %v", len(msgs), c.Size, msg)
			msgs = append(msgs, msg)
			if len(msgs) == c.Size {
				c.debugf("exceeded %d messages – flushing", c.Size)
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
			//
			// Note that this is will cause calls to the send methods to panic
			// if the client is used after being closed, definitely not ideal,
			// we should return an error like io.EOF instead. Since this has
			// been the historical behavior already I'll assume it hasn't been
			// a problem and I'll fix it later.
			close(c.msgs)
			for msg := range c.msgs {
				c.debugf("buffer (%d/%d) %v", len(msgs), c.Size, msg)
				msgs = append(msgs, msg)
			}

			c.debugf("exit requested – flushing %d", len(msgs))
			c.send(msgs)
			c.debugf("exit")
			return
		}
	}
}

// Verbose log.
func (c *Client) debugf(format string, args ...interface{}) {
	if c.Verbose {
		c.logf(format, args...)
	}
}

// Unconditional log.
func (c *Client) logf(format string, args ...interface{}) {
	c.Logger.Logf(format, args...)
}

func (c *Client) errorf(format string, args ...interface{}) {
	c.Logger.Errorf(format, args...)
}

// Return uuid string.
func uid() string {
	return uuid.NewRandom().String()
}
