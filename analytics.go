package analytics

//
// dependencies
//

import "github.com/jehiah/go-strftime"
import "github.com/twinj/uuid"
import . "encoding/json"
import "io/ioutil"
import "net/http"
import "errors"
import "bytes"
import "time"
import "log"

//
// Library version
//

const Version = "0.0.1"

//
// Default API end-point
//

const api = "https://api.segment.io"

//
// Message type.
//

type Message map[string]interface{}

//
// Segment.io client
//

type Client struct {
	Debug         bool
	BufferSize    int
	FlushInterval time.Duration
	Endpoint      string
	Key           string
	buffer        []Message
}

//
// Batch message
//

type batch struct {
	Messages  []Message `json:"batch"`
	MessageId string    `json:"messageId"`
}

//
// UUID formatting.
//

func init() {
	// TODO: wtf, this is lame
	uuid.SwitchFormat(uuid.CleanHyphen, false)
}

//
// Return a new Segment.io client
// with the given write key.
//

func New(key string) (c *Client) {
	defer func() {
		go func() {
			for {
				time.Sleep(c.FlushInterval)
				c.log("interval %v reached", c.FlushInterval)
				go c.flush()
			}
		}()
	}()

	return &Client{
		Debug:         false,
		BufferSize:    20,
		FlushInterval: 30 * time.Second,
		Key:           key,
		Endpoint:      api,
		buffer:        make([]Message, 0),
	}
}

//
// Buffer an alias message.
//

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

//
// Buffer a page message.
//

func (c *Client) Page(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "page"))

	return nil
}

//
// Buffer a screen message.
//

func (c *Client) Screen(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "screen"))

	return nil
}

//
// Buffer a group message.
//

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

//
// Buffer an identify message.
//

func (c *Client) Identify(msg Message) error {
	if msg["userId"] == nil && msg["anonymousId"] == nil {
		return errors.New("You must pass either an 'anonymousId' or 'userId'.")
	}

	c.queue(message(msg, "identify"))

	return nil
}

//
// Buffer a track message.
//

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

//
// Return a new initialized message map
// with `msg` values and context merged.
//

func message(msg Message, call string) Message {
	m := newMessage(call)

	if msg["context"] != nil {
		merge(m["context"].(map[string]interface{}), msg["context"].(map[string]interface{}))
		delete(msg, "context")
	}

	merge(m, msg)

	return m
}

//
// Return new initialzed message map.
//

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

//
// Merge two maps.
//

func merge(dst Message, src Message) {
	for k, v := range src {
		dst[k] = v
	}
}

//
// Return uuid.
//

func uid() string {
	return uuid.NewV4().String()
}

//
// Return formatted timestamp.
//

func timestamp() string {
	return strftime.Format("%Y-%m-%dT%H:%M:%S%z", time.Now())
}

//
// Log in debug mode.
//

func (c *Client) log(format string, v ...interface{}) {
	if c.Debug {
		log.Printf(format, v...)
	}
}

//
// Buffer the given message and flush
// when the buffer exceeds .BufferSize.
//

func (c *Client) queue(msg Message) {
	c.buffer = append(c.buffer, msg)

	c.log("buffer (%d/%d) %v", len(c.buffer), c.BufferSize, msg)

	if len(c.buffer) >= c.BufferSize {
		go c.flush()
	}
}

//
// Return a batch message primed
// with context properties
//

func batchMessage(msgs []Message) *batch {
	return &batch{
		MessageId: uid(),
		Messages:  msgs,
	}
}

//
// Flush the buffered messages.
//
// TODO: better error-handling,
// this is really meh, it would
// be better if we used a chan
// to deliver them.
//

func (c *Client) flush() error {
	if len(c.buffer) == 0 {
		c.log("no messages to flush")
		return nil
	}

	c.log("flushing %d messages", len(c.buffer))
	json, err := Marshal(batchMessage(c.buffer))

	if err != nil {
		c.log("error: %v", err)
		return err
	}

	c.buffer = nil

	client := &http.Client{}
	url := c.Endpoint + "/v1/import"
	c.log("POST %s with %d bytes", url, len(json))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	if err != nil {
		c.log("error: %v", err)
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(json)))
	req.SetBasicAuth(c.Key, "")

	res, err := client.Do(req)

	if err != nil {
		c.log("error: %v", err)
		return err
	}

	c.log("%d response", res.StatusCode)

	if res.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(res.Body)
		c.log("error: %s", string(body))
		c.log("error: %s", string(json))
	}

	return err
}
