package analytics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Client struct {
	*client
	config     S3ClientConfig
	apiContext *apiContext
	uploader   *s3manager.Uploader
}

// S3ClientConfig provides configuration for S3 Client.
type S3ClientConfig struct {
	Bucket string
	Stage  string

	// Stream is a name of the stream where messages will be delivered. Examples:
	// tuna, salmon, haring, etc. Each system receives its own stream.
	Stream string

	KeyConstructor func(now func() Time, uid func() string) string
}

// NewS3ClientWithConfig creates S3 client from provided configuration.
// Pass empty S3ClientConfig{} to use default config.
func NewS3ClientWithConfig(s3cfg S3ClientConfig, cfg Config) (Client, error) {
	client, err := newWithConfig("", cfg)
	if err != nil {
		return nil, err
	}

	s3Config, err := makeS3ClientConfig(s3cfg)
	if err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	c := &s3Client{
		client: client,
		config: s3Config,
		apiContext: &apiContext{
			APIID: uid(),
			Stage: s3Config.Stage,
		},
		uploader: uploader,
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

	ex := newExecutor(c.maxConcurrentRequests)
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

type apiContext struct {
	APIID        string `json:"apiId,omitempty"`
	RequestID    string `json:"requestId,omitempty"`
	ResourcePath string `json:"resourcePath,omitempty"`
	Stage        string `json:"stage,omitempty"`
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
	result := targetMessage{
		APIContext: apiContext,
		Event:      m,
		SentAt:     ts,
		ReceivedAt: ts,
	}
	type alias targetMessage
	b, err := json.Marshal(alias(result))
	if err != nil {
		return &result, err
	}
	if len(b) > maxBytes {
		return &result, ErrMessageTooBig
	}
	result.json = b
	return &result, nil
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

func (c *s3Client) flush(q *messageQueue, wg *sync.WaitGroup, ex *executor) {
	if msgs := q.flush(); msgs != nil {
		c.debugf("flushing %d messages", len(msgs))
		c.sendAsync(msgs, wg, ex)
	}
}

// Send batch request.
func (c *s3Client) send(msgs []message) {
	const attempts = 10
	var err error

	buf := &bytes.Buffer{}
	wr := gzip.NewWriter(buf)
	encoder := json.NewEncoder(wr)

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

	if err = wr.Close(); err != nil {
		c.errorf("flushing writer failed: %s", err)
		return
	}

	if buf.Len() == 0 || len(marshalledMessages) == 0 {
		c.errorf("empty buffer, send is not possible")
		return
	}

	for i := 0; i != attempts; i++ {
		if err = c.upload(buf); err == nil {
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
	c.notifyFailure(marshalledMessages, err)
}

// Upload batch to S3.
func (c *s3Client) upload(r io.Reader) error {
	key := c.config.KeyConstructor(c.now, uid)
	c.debugf("uploading to s3://%s/%s", c.config.Bucket, key)

	input := &s3manager.UploadInput{
		Body:   r,
		Bucket: &(c.config.Bucket),
		Key:    &key,
	}
	_, err := c.uploader.Upload(input)
	return err
}

func stringPtr(s string) *string { return &s }
