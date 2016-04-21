package analytics

import (
	"io"
	"io/ioutil"
	"sync"

	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jehiah/go-strftime"
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

// Message interface.
type message interface {
	setMessageId(string)
	setTimestamp(string)
}

// Message fields common to all.
type Message struct {
	Type      string `json:"type,omitempty"`
	MessageId string `json:"messageId,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	SentAt    string `json:"sentAt,omitempty"`
}

// Batch message.
type Batch struct {
	Context  map[string]interface{} `json:"context,omitempty"`
	Messages []interface{}          `json:"batch"`
	Message
}

// Identify message.
type Identify struct {
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Message
}

// Group message.
type Group struct {
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	GroupId      string                 `json:"groupId"`
	Message
}

// Track message.
type Track struct {
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Event        string                 `json:"event"`
	Message
}

// Page message.
type Page struct {
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
	Traits       map[string]interface{} `json:"properties,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Message
}

// Alias message.
type Alias struct {
	PreviousId string `json:"previousId"`
	UserId     string `json:"userId"`
	Message
}

// Client which batches messages and flushes at the given Interval or
// when the Size limit is exceeded. Set Verbose to true to enable
// logging output.
type Client struct {
	Endpoint string
	// Interval represents the duration at which messages are flushed. It may be
	// configured only before any messages are enqueued.
	Interval time.Duration
	Size     int
	Logger   Logger
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
		Logger:   newDefaultLogger(),
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

// Alias buffers an "alias" message.
func (c *Client) Alias(msg Alias) error {
	if msg.UserId == "" {
		return errors.New("You must pass a 'userId'.")
	}

	if msg.PreviousId == "" {
		return errors.New("You must pass a 'previousId'.")
	}

	msg.Type = "alias"
	c.queue(&msg)
	return nil
}

// Page buffers an "page" message.
func (c *Client) Page(msg Page) error {
	if msg.UserId == "" && msg.AnonymousId == "" {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	msg.Type = "page"
	c.queue(&msg)
	return nil
}

// Group buffers an "group" message.
func (c *Client) Group(msg Group) error {
	if msg.GroupId == "" {
		return errors.New("You must pass a 'groupId'.")
	}

	if msg.UserId == "" && msg.AnonymousId == "" {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	msg.Type = "group"
	c.queue(&msg)
	return nil
}

// Identify buffers an "identify" message.
func (c *Client) Identify(msg Identify) error {
	if msg.UserId == "" && msg.AnonymousId == "" {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	msg.Type = "identify"
	c.queue(&msg)
	return nil
}

// Track buffers an "track" message.
func (c *Client) Track(msg Track) error {
	if msg.Event == "" {
		return errors.New("You must pass 'event'.")
	}

	if msg.UserId == "" && msg.AnonymousId == "" {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	msg.Type = "track"
	c.queue(&msg)
	return nil
}

func (c *Client) startLoop() {
	go c.loop()
}

// Queue message.
func (c *Client) queue(msg message) {
	c.once.Do(c.startLoop)
	msg.setMessageId(c.uid())
	msg.setTimestamp(timestamp(c.now()))
	c.msgs <- msg
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

	batch := new(Batch)
	batch.Messages = msgs
	batch.MessageId = c.uid()
	batch.SentAt = timestamp(c.now())
	batch.Context = DefaultContext

	b, err := json.Marshal(batch)
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

	res, err := c.Client.Do(req)
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

	defer res.Body.Close()
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

// Set message timestamp if one is not already set.
func (m *Message) setTimestamp(s string) {
	if m.Timestamp == "" {
		m.Timestamp = s
	}
}

// Set message id.
func (m *Message) setMessageId(s string) {
	if m.MessageId == "" {
		m.MessageId = s
	}
}

// Return formatted timestamp.
func timestamp(t time.Time) string {
	return strftime.Format("%Y-%m-%dT%H:%M:%S%z", t)
}

// Return uuid string.
func uid() string {
	return uuid.NewRandom().String()
}
