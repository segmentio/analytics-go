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

type client struct {
	debug      bool
	key        string
	url        string
	flushAt    int
	flushAfter time.Duration
	buffer     []*interface{}
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

func Client(key string) *client {
	c := &client{
		key:    key,
		url:    api,
		buffer: make([]*interface{}, 0),
	}

	c.FlushAt(500)
	c.FlushAfter(10 * time.Second)

	return c
}

//
// Set buffer max.
//

func (c *client) FlushAt(n int) {
	c.flushAt = n
}

//
// Set buffer flush interal.
//

func (c *client) FlushAfter(interval time.Duration) {
	c.flushAfter = interval

	go func() {
		for {
			time.Sleep(interval)
			if c.debug {
				log.Printf("interval %v reached", interval)
			}
			c.flush()
		}
	}()
}

//
// Enable debug mode.
//

func (c *client) Debug() {
	c.debug = true
}

//
// Set target url
//

func (c *client) URL(url string) {
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

func (c *client) flush() error {
	if len(c.buffer) == 0 {
		if c.debug {
			log.Print("no messages to flush")
		}
		return nil
	}

	if c.debug {
		log.Printf("flushing %d messages", len(c.buffer))
	}
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

	res, err := client.Do(req)

	if res != nil {
		// TODO: how the fuck do you ignore res ^
	}

	return err
}

//
// Buffer the given message and flush
// when the buffer exceeds .flushAt.
//

func (c *client) bufferMessage(msg interface{}) error {
	c.buffer = append(c.buffer, &msg)

	if c.debug {
		log.Printf("buffer (%d/%d) %v", len(c.buffer), c.flushAt, msg)
	}

	if len(c.buffer) >= c.flushAt {
		return c.flush()
	}

	return nil
}

//
// Buffer an alias message
//

func (c *client) Alias(previousId string) error {
	return c.bufferMessage(&alias{"Alias", previousId, timestamp()})
}

//
// Buffer a page message
//

func (c *client) Page(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Page", name, category, properties, timestamp()})
}

//
// Buffer a screen message
//

func (c *client) Screen(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Screen", name, category, properties, timestamp()})
}

//
// Buffer a group message
//

func (c *client) Group(id string, traits interface{}) error {
	return c.bufferMessage(&group{"Group", id, traits, timestamp()})
}

//
// Buffer an identify message
//

func (c *client) Identify(traits interface{}) error {
	return c.bufferMessage(&identify{"Identify", traits, timestamp()})
}

//
// Buffer a track message
//

func (c *client) Track(event string, properties interface{}) error {
	return c.bufferMessage(&track{"Track", event, properties, timestamp()})
}
