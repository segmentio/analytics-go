package analytics

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type s3Client struct {
	*client
	config     S3ClientConfig
	apiContext *apiContext
	uploader   uploader
	//s3client works only with one type of msg
	tagsOnlyMsg tagsOnlyMsg
}

// S3 is a configuration for s3Client.
type S3 struct {
	Bucket string
	Stage  string

	// Stream is an arbitrary name of the stream where messages will be delivered.
	// Examples: tuna, salmon, haring, etc. Each system receives its own stream.
	Stream string

	// MaxBatchBytes size repsresents the size of buffer or file and when events are flushed
	MaxBatchBytes int

	// BufferFilePath if specified the temp file will be used to store the data
	BufferFilePath string

	KeyConstructor func(now func() Time, uid func() string) string

	UploaderOptions []func(*s3manager.Uploader)
}

// S3ClientConfig provides configuration for S3 Client.
type S3ClientConfig struct {
	Config
	S3 S3
}

// NewS3ClientWithConfig creates S3 client from provided configuration.
// Pass empty S3ClientConfig{} to use default config.
func NewS3ClientWithConfig(config S3ClientConfig) (Client, error) {
	cfg, err := makeS3ClientConfig(config)
	if err != nil {
		return nil, err
	}

	client, err := newWithConfig("", cfg.Config)
	if err != nil {
		return nil, err
	}

	client.msgs = make(chan Message, 1024*4) // overrite the buffer

	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess, cfg.S3.UploaderOptions...)

	c := &s3Client{
		client: client,
		config: cfg,
		apiContext: &apiContext{
			APIID: uid(),
			Stage: cfg.S3.Stage,
		},
		uploader: uploader,
	}
	buf, err := newBuffer(c.config.S3.BufferFilePath, c.config.S3.MaxBatchBytes)
	if err != nil {
		return nil, fmt.Errorf("can't create a buffer for the encoder: %s", err)
	}

	go c.loop(buf)     // custom implementation
	go c.loopMetrics() // reuse client's implementation

	return c, nil
}

// a copy of client.loop() function.
func (c *s3Client) loop(buf encodedBuffer) {
	defer buf.Close()

	defer close(c.shutdown)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	tick := time.NewTicker(c.Interval)
	defer tick.Stop()

	ex := newExecutor(c.maxConcurrentRequests)
	defer ex.close()

	bw := newBufferedEncoder(
		c.BatchSize,
		int64(c.config.S3.MaxBatchBytes),
		buf,
	)

	for {
		select {
		case msg := <-c.msgs:
			c.push(bw, msg, wg, ex)

		case <-tick.C:
			c.flush(bw, wg, ex)

		case <-c.quit:
			log.Println("exit requested – draining messages")
			c.debugf("exit requested – draining messages")

			// Drain the msg channel, we have to close it first so no more
			// messages can be pushed and otherwise the loop would never end.
			close(c.msgs)
			for msg := range c.msgs {
				c.push(bw, msg, wg, ex)
			}

			c.flush(bw, wg, ex)
			c.debugf("exit")
			return
		}
	}
}

func newBuffer(path string, size int) (encodedBuffer, error) {
	if path == "" {
		return newMemBuffer(size), nil
	}

	return newFileBuffer(path)
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
}

func (m *targetMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *targetMessage) Msg() Message {
	return m.Event
}

func (m *targetMessage) size() int {
	return -1
}

// dummy message to store flags
type tagsOnlyMsg struct {
	t []string
}

func (m *tagsOnlyMsg) tags() []string {
	return m.t
}

func (m *tagsOnlyMsg) validate() error {
	return nil
}

func (c *s3Client) push(encoder *bufferedEncoder, m Message, wg *sync.WaitGroup, ex *executor) {
	c.setTagsIfExsist(m)

	ready, err := encodeMessage(encoder, m, c.apiContext, c.now)
	if err != nil {
		c.errorf("cant encode message: ", err)
		c.notifyFailureMsg(m, err, 1)
	}
	c.debugf("buffer (%d/%d) %v", encoder.messages, c.BatchSize, m)

	if ready {
		c.debugf("exceeded messages batch limit with batch of %d messages – flushing", encoder.messages)
		c.send(encoder)
	}
}

func (c *s3Client) setTagsIfExsist(m Message) {
	if len(c.tagsOnlyMsg.t) == 0 {
		c.tagsOnlyMsg.t = m.tags()
	}
}

func encodeMessage(bw *bufferedEncoder, m Message, ctx *apiContext, now func() Time) (ready bool, err error) {
	ts := now()
	msg := targetMessage{
		APIContext: ctx,
		Event:      m,
		SentAt:     ts,
		ReceivedAt: ts,
	}
	type alias targetMessage // we won't use json.Marshaller implementation

	return bw.Push(alias(msg))
}

// Asychronously send a batched requests.
func (c *s3Client) sendAsync(bw *bufferedEncoder, wg *sync.WaitGroup, ex *executor) {
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
		c.send(bw)
	}) {
		wg.Done()
		c.errorf("sending messages failed - %s", ErrTooManyRequests)
		c.notifyFailureMsg(&c.tagsOnlyMsg, ErrTooManyRequests, int64(bw.TotalMsgs()))
	}
}

func (c *s3Client) flush(bw *bufferedEncoder, wg *sync.WaitGroup, ex *executor) {
	msgs := bw.TotalMsgs()
	if msgs > 0 {
		c.debugf("flushing %d messages", msgs)
		c.send(bw)
	}
}

// Send batch request.
func (c *s3Client) send(bw *bufferedEncoder) {
	const attempts = 10

	if bw.BytesLen() == 0 {
		c.errorf("empty buffer, send is not possible")
		return
	}

	defer bw.Reset()

	msgs := bw.TotalMsgs()
	for i := 0; i != attempts; i++ {
		reader, err := bw.Reader()
		if err != nil {
			c.errorf("can't get reader", err)
		}

		if err = c.upload(reader); err == nil {
			c.notifySuccessMsg(&c.tagsOnlyMsg, int64(msgs))
			return
		}

		// Wait for either a retry timeout or the client to be closed.
		select {
		case <-time.After(c.RetryAfter(i)):
			c.errorf("%d messages dropped because of error: %s", msgs, err)
			return
		case <-c.quit:
			c.errorf("%d messages dropped because they failed to be sent and the client was closed, upload error: %s", msgs, err)
			return
		}
	}

	c.errorf("%d messages dropped because they failed to be sent after %d attempts", msgs, attempts)
}

// Upload batch to S3.
func (c *s3Client) upload(r io.Reader) error {
	key := c.config.S3.KeyConstructor(c.now, uid)
	c.debugf("uploading to s3://%s/%s", c.config.S3.Bucket, key)

	input := &s3manager.UploadInput{
		Body:   r,
		Bucket: aws.String(c.config.S3.Bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(key),
	}
	_, err := c.uploader.Upload(input)
	return err
}

type bufferedEncoder struct {
	maxBatchSize  int
	maxBatchBytes int64

	buf      encodedBuffer
	encoder  *json.Encoder
	gziper   *gzip.Writer
	messages int
}

func newBufferedEncoder(maxBatchSize int, maxBatchBytes int64, buf encodedBuffer) *bufferedEncoder {
	w := &bufferedEncoder{
		maxBatchSize:  maxBatchSize,
		maxBatchBytes: maxBatchBytes,
		buf:           buf,
	}

	w.gziper = gzip.NewWriter(w.buf)
	w.encoder = json.NewEncoder(w.gziper)
	return w
}

func (q *bufferedEncoder) BytesLen() int64 {
	return q.buf.Size()
}

func (q *bufferedEncoder) TotalMsgs() int {
	return q.messages
}

func (q *bufferedEncoder) Reader() (io.Reader, error) {
	err := q.gziper.Close()
	if err != nil {
		return nil, err
	}
	return q.buf.Reader()
}

func (q *bufferedEncoder) Push(v interface{}) (bool, error) {
	if err := q.encoder.Encode(v); err != nil {
		return false, err
	}
	q.messages++

	if q.buf.Size() >= q.maxBatchBytes {
		return true, nil
	}

	if q.messages >= q.maxBatchSize {
		return true, nil
	}

	return false, nil
}

func (q *bufferedEncoder) Reset() error {
	err := q.buf.Reset()
	if err != nil {
		return err
	}
	q.gziper = gzip.NewWriter(q.buf)
	q.encoder = json.NewEncoder(q.gziper)
	q.messages = 0
	return nil
}

type encodedBuffer interface {
	io.WriteCloser
	Size() int64
	Reader() (io.Reader, error)
	Reset() error
}

type memBuffer struct {
	buf *bytes.Buffer
}

func newMemBuffer(size int) *memBuffer {
	var buf bytes.Buffer
	if size > 1 {
		buf.Grow(size)
	}

	return &memBuffer{
		buf: &buf,
	}
}

func (m *memBuffer) Write(p []byte) (n int, err error) {
	return m.buf.Write(p)
}

func (m *memBuffer) Reader() (io.Reader, error) {
	return bytes.NewReader(m.buf.Bytes()), nil
}

func (m *memBuffer) Reset() error {
	m.buf.Reset()
	return nil
}

func (m *memBuffer) Size() int64 {
	return int64(m.buf.Len())
}

func (m *memBuffer) Close() error {
	return nil
}

type fileBuffer struct {
	fd     *os.File
	writer *bufio.Writer
	reader *bufio.Reader
	size   int64
}

func newFileBuffer(path string) (*fileBuffer, error) {
	fd, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &fileBuffer{
		fd:     fd,
		writer: bufio.NewWriter(fd),
		reader: bufio.NewReader(fd),
	}, nil
}

func (m *fileBuffer) Write(p []byte) (n int, err error) {
	n, err = m.writer.Write(p)
	if err != nil {
		return n, err
	}

	m.size += int64(n)
	return n, nil
}

func (m *fileBuffer) Reader() (io.Reader, error) {
	if err := m.writer.Flush(); err != nil {
		return nil, err
	}

	if _, err := m.fd.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return io.LimitReader(m.reader, m.size), nil
}

func (m *fileBuffer) Reset() error {
	m.size = 0
	if _, err := m.fd.Seek(0, io.SeekStart); err != nil {
		return err
	}
	m.writer.Reset(m.fd)
	m.reader.Reset(m.fd)

	return nil
}

func (m *fileBuffer) Size() int64 {
	return m.size
}

func (m *fileBuffer) Close() error {
	return m.fd.Close()
}
