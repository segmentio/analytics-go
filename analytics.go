package analytics

//
// dependencies
//

import "github.com/jehiah/go-strftime"
import "github.com/nu7hatch/gouuid"
import . "encoding/json"
import "net/http"
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
// Segment.io client
//

type Client struct {
	Debug         bool
	BufferSize    int
	running       bool
	key           string
	url           string
	flushInterval time.Duration
	buffer        []*interface{}
}

//
// Message context library
//

type contextLibrary struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//
// Message context
//

type context struct {
	Library contextLibrary `json:"library"`
}

//
// Identify message
//

type identify struct {
	Action    string      `json:"action"`
	Traits    interface{} `json:"trailts"`
	Timestamp string      `json:"timestamp"`
}

//
// Alias message
//

type alias struct {
	Action     string `json:"action"`
	PreviousId string `json:"previousId"`
	Timestamp  string `json:"timestamp"`
}

//
// Track message
//

type track struct {
	Action     string      `json:"action"`
	Event      string      `json:"event"`
	Properties interface{} `json:"properties"`
	Timestamp  string      `json:"timestamp"`
}

//
// Group message
//

type group struct {
	Action    string      `json:"action"`
	GroupId   string      `json:"groupId"`
	Traits    interface{} `json:"trailts"`
	Timestamp string      `json:"timestamp"`
}

//
// Page message
//

type page struct {
	Action     string      `json:"action"`
	Category   string      `json:"category"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
	Timestamp  string      `json:"timestamp"`
}

//
// Batch message
//

type batch struct {
	Context   context        `json:"context"`
	RequestId string         `json:"requestId"`
	Messages  []*interface{} `json:"batch"`
}

//
// Return a new Segment.io client
// with the given write key.
//

func New(key string) *Client {
	c := &Client{
		Debug:      false,
		BufferSize: 500,
		key:        key,
		url:        api,
		buffer:     make([]*interface{}, 0),
	}

	c.FlushAfter(10 * time.Second)

	return c
}

//
// Set buffer flush interal.
//

func (c *Client) FlushAfter(interval time.Duration) {
	c.flushInterval = interval

	if c.running {
		return
	}

	c.running = true

	go func() {
		for {
			time.Sleep(c.flushInterval)
			c.log("interval %v reached", c.flushInterval)
			c.flush()
		}
	}()
}

//
// Set target url
//

func (c *Client) URL(url string) {
	c.url = url
}

//
// Return formatted timestamp.
//

func timestamp() string {
	return strftime.Format("%Y-%m-%dT%H:%M:%S%z", time.Now())
}

// Return a batch message primed
// with context properties
//

func createBatch(msgs []*interface{}) (*batch, error) {
	uid, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	batch := &batch{
		RequestId: uid.String(),
		Messages:  msgs,
		Context: context{
			Library: contextLibrary{
				Name:    "analytics-go",
				Version: Version,
			},
		},
	}

	return batch, nil
}

//
// Flush the buffered messages.
//

func (c *Client) flush() error {
	if len(c.buffer) == 0 {
		c.log("no messages to flush")
		return nil
	}

	c.log("flushing %d messages", len(c.buffer))
	batch, err := createBatch(c.buffer)

	if err != nil {
		return err
	}

	json, err := Marshal(batch)

	if err != nil {
		return err
	}

	c.buffer = nil

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.url+"/v1/batch", bytes.NewBuffer(json))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(json)))
	req.SetBasicAuth(c.key, "")

	_, err = client.Do(req)

	return err
}

//
// Buffer the given message and flush
// when the buffer exceeds .BufferSize.
//

func (c *Client) bufferMessage(msg interface{}) error {
	c.buffer = append(c.buffer, &msg)

	c.log("buffer (%d/%d) %v", len(c.buffer), c.BufferSize, msg)

	if len(c.buffer) >= c.BufferSize {
		return c.flush()
	}

	return nil
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
// Buffer an alias message
//

func (c *Client) Alias(previousId string) error {
	return c.bufferMessage(&alias{"Alias", previousId, timestamp()})
}

//
// Buffer a page message
//

func (c *Client) Page(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Page", name, category, properties, timestamp()})
}

//
// Buffer a screen message
//

func (c *Client) Screen(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Screen", name, category, properties, timestamp()})
}

//
// Buffer a group message
//

func (c *Client) Group(id string, traits interface{}) error {
	return c.bufferMessage(&group{"Group", id, traits, timestamp()})
}

//
// Buffer an identify message
//

func (c *Client) Identify(traits interface{}) error {
	return c.bufferMessage(&identify{"Identify", traits, timestamp()})
}

//
// Buffer a track message
//

func (c *Client) Track(event string, properties interface{}) error {
	return c.bufferMessage(&track{"Track", event, properties, timestamp()})
}
