package analytics

//
// dependencies
//

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
	Action string      `json:"action"`
	Traits interface{} `json:"trailts"`
}

//
// Alias message
//

type alias struct {
	Action     string `json:"action"`
	PreviousId string `json:"previousId"`
}

//
// Track message
//

type track struct {
	Action     string      `json:"action"`
	Event      string      `json:"event"`
	Properties interface{} `json:"properties"`
}

//
// Group message
//

type group struct {
	Action  string      `json:"action"`
	GroupId string      `json:"groupId"`
	Traits  interface{} `json:"trailts"`
}

//
// Page message
//

type page struct {
	Action     string      `json:"action"`
	Category   string      `json:"category"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}

//
// Batch message
//

type batch struct {
	Context   context        `json:"context"`
	RequestId string         `json:"requestId"`
	Messages  []*interface{} `json:"messages"`
}

//
// Return a new Segment.io client
// with the given write key.
//

func Client(key string) *client {
	return &client{
		key:        key,
		url:        api,
		flushAt:    500,
		flushAfter: 10 * time.Second,
		buffer:     make([]*interface{}, 0),
	}
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

// Return a batch message primed
// with context properties
//

func createBatch(msgs []*interface{}) (*batch, error) {
	uid, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	batch := &batch{
		Context: context{
			Library: contextLibrary{
				Name:    "analytics-go",
				Version: Version,
			},
		},
		RequestId: uid.String(),
		Messages:  msgs,
	}

	return batch, nil
}

//
// Flush the buffered messages.
//

func (c *client) flush() error {
	b, err := createBatch(c.buffer)

	if err != nil {
		return err
	}

	j, err := Marshal(b)

	if err != nil {
		return err
	}

	c.buffer = nil

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.url+"/v1/batch", bytes.NewBuffer(j))

	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "analytics-go (version: "+Version+")")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", string(len(j)))
	req.SetBasicAuth(c.key, "")

	res, err := client.Do(req)

	if res != nil {
		// TODO: how the fuck ^
	}

	return nil
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
		log.Printf("flushing %d messages", len(c.buffer))
		return c.flush()
	}

	return nil
}

//
// Buffer an alias message
//

func (c *client) Alias(previousId string) error {
	return c.bufferMessage(&alias{"Alias", previousId})
}

//
// Buffer a page message
//

func (c *client) Page(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Page", name, category, properties})
}

//
// Buffer a screen message
//

func (c *client) Screen(name string, category string, properties interface{}) error {
	return c.bufferMessage(&page{"Screen", name, category, properties})
}

//
// Buffer a group message
//

func (c *client) Group(id string, traits interface{}) error {
	return c.bufferMessage(&group{"Group", id, traits})
}

//
// Buffer an identify message
//

func (c *client) Identify(traits interface{}) error {
	return c.bufferMessage(&identify{"Identify", traits})
}

//
// Buffer a track message
//

func (c *client) Track(event string, properties interface{}) error {
	// TODO: .timestamp ISO-8601-formatted string.
	return c.bufferMessage(&track{"Track", event, properties})
}
