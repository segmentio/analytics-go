package analytics

//
// dependencies
//

import "github.com/jehiah/go-strftime"
import "github.com/twinj/uuid"
import . "encoding/json"
import "io/ioutil"
import "net/http"
import "reflect"
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
	Type        string      `json:"type"`
	Traits      interface{} `json:"trailts"`
	Timestamp   string      `json:"timestamp"`
	UserId      string      `json:"userId"`
	AnonymousId string      `json:"anonymousId"`
	Context     *context    `json:"context"`
	MessageId   string      `json:"messageId"`
}

//
// Alias message
//

type alias struct {
	Type       string   `json:"type"`
	PreviousId string   `json:"previousId"`
	Timestamp  string   `json:"timestamp"`
	Context    *context `json:"context"`
	MessageId  string   `json:"messageId"`
}

//
// Track message
//

type track struct {
	Type        string      `json:"type"`
	Event       string      `json:"event"`
	Properties  interface{} `json:"properties"`
	Timestamp   string      `json:"timestamp"`
	UserId      string      `json:"userId"`
	AnonymousId string      `json:"anonymousId"`
	Context     *context    `json:"context"`
	MessageId   string      `json:"messageId"`
}

//
// Group message
//

type group struct {
	Type        string      `json:"type"`
	GroupId     string      `json:"groupId"`
	Traits      interface{} `json:"trailts"`
	Timestamp   string      `json:"timestamp"`
	UserId      string      `json:"userId"`
	AnonymousId string      `json:"anonymousId"`
	Context     *context    `json:"context"`
	MessageId   string      `json:"messageId"`
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
	Context    *context    `json:"context"`
	MessageId  string      `json:"messageId"`
}

//
// Batch message
//

type batch struct {
	Messages  []*interface{} `json:"batch"`
	MessageId string         `json:"messageId"`
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
	ctx := messageContext()
	c.bufferMessage(&alias{"alias", previousId, timestamp(), ctx, uid()})
}

//
// Buffer a page message
//

func (c *Client) Page(name string, category string, properties interface{}) {
	ctx := messageContext()
	c.bufferMessage(&page{"page", name, category, properties, timestamp(), ctx, uid()})
}

//
// Buffer a screen message
//

func (c *Client) Screen(name string, category string, properties interface{}) {
	ctx := messageContext()
	c.bufferMessage(&page{"screen", name, category, properties, timestamp(), ctx, uid()})
}

//
// Buffer a group message
//

func (c *Client) Group(id string, traits interface{}) {
	user, anon := ids(traits)
	ctx := messageContext()
	c.bufferMessage(&group{"group", id, traits, timestamp(), user, anon, ctx, uid()})
}

//
// Buffer an identify message
//

func (c *Client) Identify(traits interface{}) {
	user, anon := ids(traits)
	ctx := messageContext()
	c.bufferMessage(&identify{"identify", traits, timestamp(), user, anon, ctx, uid()})
}

//
// Buffer a track message
//

func (c *Client) Track(event string, properties interface{}) {
	user, anon := ids(properties)
	ctx := messageContext()
	c.bufferMessage(&track{"track", event, properties, timestamp(), user, anon, ctx, uid()})
}

//
// Return uuid.
//

func uid() string {
	return uuid.NewV4().String()
}

//
// Return UserId or AnonymousId field value.
//

func ids(properties interface{}) (string, string) {
	userId := ""
	anonId := ""

	val := reflect.ValueOf(properties)

	user := val.FieldByName("UserId")

	if user.IsValid() {
		userId = user.String()
	}

	anon := val.FieldByName("AnonymousId")

	if anon.IsValid() {
		anonId = anon.String()
	}

	return userId, anonId
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

//
// Return message context.
//

func messageContext() *context {
	return &context{
		Library: contextLibrary{
			Name:    "analytics-go",
			Version: Version,
		},
	}
}

//
// Return a batch message primed
// with context properties
//

func batchMessage(msgs []*interface{}) *batch {
	return &batch{
		MessageId: uid(),
		Messages:  msgs,
	}
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
	json, err := Marshal(batchMessage(c.buffer))

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
		c.log("error: %s", string(json))
	}

	return err
}
