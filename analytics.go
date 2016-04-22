package analytics

import (
	"io/ioutil"
	"os"
	"sync"

	"bytes"
	"encoding/json"
	"log"
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
	Interval time.Duration
	Size     int
	Logger   *log.Logger
	Verbose  bool
	Client   http.Client
	key      string
	msgs     chan interface{}
	quit     chan struct{}
	shutdown chan struct{}
	uid      func() string
	now      func() time.Time
	once     sync.Once
}

// New client with write key.
func New(key string) *Client {
	c := &Client{
		Endpoint: Endpoint,
		Interval: 5 * time.Second,
		Size:     250,
		Logger:   log.New(os.Stderr, "segment ", log.LstdFlags),
		Verbose:  false,
		Client:   *http.DefaultClient,
		key:      key,
		msgs:     make(chan interface{}, 100),
		quit:     make(chan struct{}),
		shutdown: make(chan struct{}),
		now:      time.Now,
		uid:      uid,
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
func (c *Client) Close() error {
	c.quit <- struct{}{}
	close(c.msgs)
	<-c.shutdown
	return nil
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
		c.logf("error marshalling msgs: %s", err)
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
		c.logf("error creating request: %s", err)
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(b)))
	req.SetBasicAuth(c.key, "")

	res, err := c.Client.Do(req)
	if err != nil {
		c.logf("error sending request: %s", err)
		return err
	}
	defer res.Body.Close()

	c.report(res)

	return nil
}

// Report on response body.
func (c *Client) report(res *http.Response) {
	if res.StatusCode < 400 {
		c.verbose("response %s", res.Status)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.logf("error reading response body: %s", err)
		return
	}

	c.logf("response %s: %d – %s", res.Status, res.StatusCode, body)
}

// Batch loop.
func (c *Client) loop() {
	var msgs []interface{}
	tick := time.NewTicker(c.Interval)

	for {
		select {
		case msg := <-c.msgs:
			c.verbose("buffer (%d/%d) %v", len(msgs), c.Size, msg)
			msgs = append(msgs, msg)
			if len(msgs) == c.Size {
				c.verbose("exceeded %d messages – flushing", c.Size)
				c.send(msgs)
				msgs = nil
			}
		case <-tick.C:
			if len(msgs) > 0 {
				c.verbose("interval reached - flushing %d", len(msgs))
				c.send(msgs)
				msgs = nil
			} else {
				c.verbose("interval reached – nothing to send")
			}
		case <-c.quit:
			tick.Stop()
			c.verbose("exit requested – draining msgs")
			// drain the msg channel.
			for msg := range c.msgs {
				c.verbose("buffer (%d/%d) %v", len(msgs), c.Size, msg)
				msgs = append(msgs, msg)
			}
			c.verbose("exit requested – flushing %d", len(msgs))
			c.send(msgs)
			c.verbose("exit")
			c.shutdown <- struct{}{}
			return
		}
	}
}

// Verbose log.
func (c *Client) verbose(msg string, args ...interface{}) {
	if c.Verbose {
		c.Logger.Printf(msg, args...)
	}
}

// Unconditional log.
func (c *Client) logf(msg string, args ...interface{}) {
	c.Logger.Printf(msg, args...)
}

// Return uuid string.
func uid() string {
	return uuid.NewRandom().String()
}
