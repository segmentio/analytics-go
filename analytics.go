package analytics

//
// dependencies
//

import "github.com/jehiah/go-strftime"
import "github.com/nu7hatch/gouuid"
import . "encoding/json"
import "io/ioutil"
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
	FlushInterval time.Duration
	Endpoint      string
	Key           string
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
	Type      string      `json:"type"`
	Traits    interface{} `json:"trailts"`
	Timestamp string      `json:"timestamp"`
}

//
// Alias message
//

type alias struct {
	Type       string `json:"type"`
	PreviousId string `json:"previousId"`
	Timestamp  string `json:"timestamp"`
}

//
// Track message
//

type track struct {
	Type       string      `json:"type"`
	Event      string      `json:"event"`
	Properties interface{} `json:"properties"`
	Timestamp  string      `json:"timestamp"`
}

//
// Group message
//

type group struct {
	Type      string      `json:"type"`
	GroupId   string      `json:"groupId"`
	Traits    interface{} `json:"trailts"`
	Timestamp string      `json:"timestamp"`
}

//
// Page message
//

type page struct {
	Type       string      `json:"type"`
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
		BufferSize:    500,
		FlushInterval: 30 * time.Second,
		Key:           key,
		Endpoint:      api,
		buffer:        make([]*interface{}, 0),
	}
}

//
// Buffer an alias message
//

func (c *Client) Alias(previousId string) {
	c.bufferMessage(&alias{"alias", previousId, timestamp()})
}

//
// Buffer a page message
//

func (c *Client) Page(name string, category string, properties interface{}) {
	c.bufferMessage(&page{"page", name, category, properties, timestamp()})
}

//
// Buffer a screen message
//

func (c *Client) Screen(name string, category string, properties interface{}) {
	c.bufferMessage(&page{"screen", name, category, properties, timestamp()})
}

//
// Buffer a group message
//

func (c *Client) Group(id string, traits interface{}) {
	c.bufferMessage(&group{"group", id, traits, timestamp()})
}

//
// Buffer an identify message
//

func (c *Client) Identify(traits interface{}) {
	c.bufferMessage(&identify{"identify", traits, timestamp()})
}

//
// Buffer a track message
//

func (c *Client) Track(event string, properties interface{}) {
	c.bufferMessage(&track{"track", event, properties, timestamp()})
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

func (c *Client) bufferMessage(msg interface{}) {
	c.buffer = append(c.buffer, &msg)

	c.log("buffer (%d/%d) %v", len(c.buffer), c.BufferSize, msg)

	if len(c.buffer) >= c.BufferSize {
		go c.flush()
	}
}

// Return a batch message primed
// with context properties
//

func batchMessage(msgs []*interface{}) (*batch, error) {
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
	batch, err := batchMessage(c.buffer)

	if err != nil {
		return err
	}

	json, err := Marshal(batch)

	if err != nil {
		return err
	}

	c.buffer = nil

	client := &http.Client{}
	url := c.Endpoint + "/v1/import"
	c.log("POST %s with %d bytes", url, len(json))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(json)))
	req.SetBasicAuth(c.Key, "")

	res, err := client.Do(req)
	c.log("%d response", res.StatusCode)

	if res.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(res.Body)
		c.log("error: %s", string(body))
	}

	return err
}
