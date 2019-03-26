package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type s3Client struct {
	*client
	apiContext *apiContext
}

// S3ClientConfig provides configuration for S3 Client.
type S3ClientConfig struct {
}

// NewS3ClientWithConfig creates S3 client from provided configuration.
// Pass empty S3ClientConfig{} to use default config.
func NewS3ClientWithConfig(writeKey string, config S3ClientConfig) (Client, error) {
	client, err := newWithConfig(writeKey, Config{})
	if err != nil {
		return nil, err
	}

	c := &s3Client{
		client: client,
		apiContext: &apiContext{
			Identity: identity{
				APIKey: writeKey,
			},
		},
	}

	go c.loop()        // custom implementation
	go c.loopMetrics() // reuse client's implementation

	return c, nil
}

// a copy of client.loop() function.
func (c *s3Client) loop() {
	defer close(c.shutdown)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	tick := time.NewTicker(c.Interval)
	defer tick.Stop()

	ex := newExecutor(1)
	defer ex.close()

	mq := messageQueue{
		maxBatchSize:  c.BatchSize,
		maxBatchBytes: c.maxBatchBytes(),
	}

	for {
		select {
		case msg := <-c.msgs:
			c.push(&mq, msg, wg, ex)

		case <-tick.C:
			c.flush(&mq, wg, ex)

		case <-c.quit:
			c.debugf("exit requested – draining messages")

			// Drain the msg channel, we have to close it first so no more
			// messages can be pushed and otherwise the loop would never end.
			close(c.msgs)
			for msg := range c.msgs {
				c.push(&mq, msg, wg, ex)
			}

			c.flush(&mq, wg, ex)
			c.debugf("exit")
			return
		}
	}
}

func (c *s3Client) push(q *messageQueue, m Message, wg *sync.WaitGroup, ex *executor) {
	var msg message
	var err error

	if msg, err = makeTargetMessage(m, maxMessageBytes, c.apiContext, c.now); err != nil {
		c.errorf("%s - %v", err, m)
		c.notifyFailure([]message{msg}, err)
		return
	}

	c.debugf("buffer (%d/%d) %v", len(q.pending), c.BatchSize, m)

	if msgs := q.push(msg); msgs != nil {
		c.debugf("exceeded messages batch limit with batch of %d messages – flushing", len(msgs))
		c.sendAsync(msgs, wg, ex)
	}
}

type identity struct {
	APIKey    string `json:"apiKey,omitempty"`
	Country   string `json:"country,omitempty"`
	IsDesktop *bool  `json:"isDesktop,omitempty"`
	IsMobile  *bool  `json:"isMobile,omitempty"`
	IsTablet  *bool  `json:"isTablet,omitempty"`
}
type apiContext struct {
	APIID        string   `json:"apiId,omitempty"`
	HTTPMethod   string   `json:"httpMethod,omitempty"`
	Identity     identity `json:"identity,omitempty"`
	RequestID    string   `json:"requestId,omitempty"`
	ResourceID   string   `json:"resourceId,omitempty"`
	ResourceMeta string   `json:"resourceMeta,omitempty"`
	ResourcePath string   `json:"resourcePath,omitempty"`
	Stage        string   `json:"stage,omitempty"`
}

// targetMessage is a single non-batched message delivered to s3 in one row of json.
type targetMessage struct {
	APIContext *apiContext `json:"context,omitempty"`
	Event      Message     `json:"event"`
	SentAt     Time        `json:"sentAt"`
	ReceivedAt Time        `json:"receivedAt"`

	json []byte
}

func (m *targetMessage) MarshalJSON() ([]byte, error) {
	return m.json, nil
}

func (m *targetMessage) Msg() Message {
	return m.Event
}

func (m *targetMessage) size() int {
	return len(m.json)
}

// makeTargetMessage constructs targetMessage instance.
func makeTargetMessage(m Message, maxBytes int, apiContext *apiContext, now func() Time) (message, error) {
	ts := now()
	result := &targetMessage{
		APIContext: apiContext,
		Event:      m,
		SentAt:     ts,
		ReceivedAt: ts,
	}
	b, err := json.Marshal(struct{ *targetMessage }{result})
	if err != nil {
		return result, err
	}
	if len(b) > maxBytes {
		return result, ErrMessageTooBig
	}
	result.json = b
	return result, nil
}

// Asychronously send a batched requests.
func (c *s3Client) sendAsync(msgs []message, wg *sync.WaitGroup, ex *executor) {
	wg.Add(1)

	if !ex.do(func() {
		defer wg.Done()
		defer func() {
			// In case a bug is introduced in the send function that triggers
			// a panic, we don't want this to ever crash the application so we
			// catch it here and log it instead.
			if err := recover(); err != nil {
				c.errorf("panic - %s", err)
			}
		}()
		c.send(msgs)
	}) {
		wg.Done()
		c.errorf("sending messages failed - %s", ErrTooManyRequests)
		c.notifyFailure(msgs, ErrTooManyRequests)
	}
}

// Send batch request.
func (c *s3Client) send(msgs []message) {
	const attempts = 10
	var err error

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)

	marshalledMessages := []message{}
	failedMessages := []message{}
	var lastError error

	for _, m := range msgs {
		err = encoder.Encode(m)
		if err != nil {
			failedMessages = append(failedMessages, m)
			lastError = err
		} else {
			marshalledMessages = append(marshalledMessages, m)
		}
	}
	if len(failedMessages) > 0 {
		c.errorf("marshalling message - %s", lastError)
		c.notifyFailure(failedMessages, lastError)
	}

	if buf.Len() == 0 || len(marshalledMessages) == 0 {
		c.errorf("empty buffer, send is not possible")
		return
	}

	b := buf.Bytes()

	for i := 0; i != attempts; i++ {
		if err = c.upload(b); err == nil {
			c.notifySuccess(marshalledMessages)
			return
		}

		// Wait for either a retry timeout or the client to be closed.
		select {
		case <-time.After(c.RetryAfter(i)):
		case <-c.quit:
			err = fmt.Errorf("%d messages dropped because they failed to be sent and the client was closed", len(msgs))
			c.errorf(err.Error())
			c.notifyFailure(marshalledMessages, err)
			return
		}
	}

	c.errorf("%d messages dropped because they failed to be sent after %d attempts", len(msgs), attempts)
	c.notifyFailure(msgs, err)
}

// Upload batch to S3.
func (c *s3Client) upload(b []byte) error {
	return fmt.Errorf("not implemented")
}
