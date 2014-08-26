package analytics

import . "github.com/visionmedia/go-debug"
import "github.com/jehiah/go-strftime"
import "github.com/xtgo/uuid"
import . "encoding/json"
import "io/ioutil"
import "net/http"
import "errors"
import "bytes"
import "sync"
import "time"

// Library version.
const Version = "0.0.2"

// Default API end-point.
const api = "https://api.segment.io"

// Message type.
type Message map[string]interface{}

// Debug.
var debug = Debug("analytics")

// Segment.io client
type Client struct {
	FlushAt    int
	FlushAfter time.Duration
	Endpoint   string
	Key        string
	buffer     []Message
	sync.Mutex
}

// Batch message.
type batch struct {
	Messages  []Message `json:"batch"`
	MessageId string    `json:"messageId"`
}

// New creates a new Segment.io client
// with the given write key.
func New(key string) *Client {
	c := &Client{
		FlushAt:    20,
		FlushAfter: 5 * time.Second,
		Key:        key,
		Endpoint:   api,
		buffer:     make([]Message, 0),
	}

	go c.Start()

	return c
}

// Start flusher.
func (c *Client) Start() {
	go func() {
		for {
			time.Sleep(c.FlushAfter)
			debug("interval %v reached", c.FlushAfter)
			c.flush()
		}
	}()
}

// Alias buffers an "alias" message.
func (c *Client) Alias(msg Message) error {
	if msg["userId"] == nil {
		return errors.New("You must pass a 'userId'.")
	}

	if msg["previousId"] == nil {
		return errors.New("You must pass a 'previousId'.")
	}

	c.queue(message(msg, "alias"))

	return nil
}

// Page buffers an "page" message.
func (c *Client) Page(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "page"))

	return nil
}

// Screen buffers an "screen" message.
func (c *Client) Screen(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "screen"))

	return nil
}

// Group buffers an "group" message.
func (c *Client) Group(msg Message) error {
	if msg["groupId"] == nil {
		return errors.New("You must pass a 'groupId'.")
	}

	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "group"))

	return nil
}

// Identify buffers an "identify" message.
func (c *Client) Identify(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "identify"))

	return nil
}

// Track buffers an "track" message.
func (c *Client) Track(msg Message) error {
	if msg["event"] == nil {
		return errors.New("You must pass 'event'.")
	}

	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "track"))

	return nil
}

// Return a new initialized message map
// with `msg` values and context merged.
func message(msg Message, call string) Message {
	m := newMessage(call)

	if msg["context"] != nil {
		merge(m["context"].(map[string]interface{}), msg["context"].(map[string]interface{}))
		delete(msg, "context")
	}

	merge(m, msg)

	return m
}

// Return new initialzed message map.
func newMessage(call string) Message {
	return Message{
		"type":      call,
		"timestamp": timestamp(),
		"messageId": uid(),
		"context": map[string]interface{}{
			"version": Version,
			"library": "analytics-go",
		},
	}
}

// Merge two maps.
func merge(dst Message, src Message) {
	for k, v := range src {
		dst[k] = v
	}
}

// Return uuid.
func uid() string {
	return uuid.NewRandom().String()
}

// Return formatted timestamp.
func timestamp() string {
	return strftime.Format("%Y-%m-%dT%H:%M:%S%z", time.Now())
}

// Buffer the given message and flush
// when the buffer exceeds .FlushAt.
func (c *Client) queue(msg Message) {
	c.Lock()
	defer c.Unlock()

	c.buffer = append(c.buffer, msg)

	debug("buffer (%d/%d) %v", len(c.buffer), c.FlushAt, msg)

	if len(c.buffer) >= c.FlushAt {
		go c.flush()
	}
}

// Return a batch message primed
// with context properties.
func batchMessage(msgs []Message) *batch {
	return &batch{
		MessageId: uid(),
		Messages:  msgs,
	}
}

// Flush the buffered messages.
//
// TODO: better error-handling,
// this is really meh, it would
// be better if we used a chan
// to deliver them.
//
func (c *Client) flush() error {
	c.Lock()

	if len(c.buffer) == 0 {
		debug("no messages to flush")
		c.Unlock()
		return nil
	}

	debug("flushing %d messages", len(c.buffer))
	json, err := Marshal(batchMessage(c.buffer))

	if err != nil {
		debug("error: %v", err)
		c.Unlock()
		return err
	}

	c.buffer = nil
	c.Unlock()

	client := &http.Client{}
	url := c.Endpoint + "/v1/import"
	debug("POST %s with %d bytes", url, len(json))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	if err != nil {
		debug("error: %v", err)
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(json)))
	req.SetBasicAuth(c.Key, "")

	res, err := client.Do(req)
	if err != nil {
		debug("error: %v", err)
		return err
	}
	defer res.Body.Close()

	debug("%d response", res.StatusCode)

	if res.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(res.Body)
		debug("error: %s", string(body))
		debug("error: %s", string(json))
	}

	return err
}
