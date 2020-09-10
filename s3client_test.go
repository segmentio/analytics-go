package analytics

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTargetMessageMarshalling(t *testing.T) {
	type msg struct {
		Index int `json:"index"`
		Qwer  int `json:"qwer"`
	}
	m := TrackObj{
		Track: Track{
			Event:  "FooBared",
			UserId: "tuna",
		},
		Properties: msg{
			Index: 1,
			Qwer:  3424,
		},
	}

	buf := &memBuffer{
		buf: bytes.NewBuffer(nil),
	}

	encoder := newBufferedEncoder(100, 10, buf)

	_, err := encodeMessage(encoder, m, nil, func() Time { return Time{} })
	if err != nil {
		t.Error(err)
	}

	reader, err := encoder.Reader()
	if err != nil {
		t.Error(err)
	}

	result := readAndUngzip(t, reader)

	expected := `{"event":{"userId":"tuna","event":"FooBared","timestamp":0,"properties":{"index":1,"qwer":3424}},"sentAt":0,"receivedAt":0}` + "\n"

	require.Equal(t, expected, string(result))
}

func readAndUngzip(t *testing.T, r io.Reader) []byte {
	q, err := gzip.NewReader(r)
	assert.NoError(t, err)
	defer q.Close()

	d, err := ioutil.ReadAll(q)
	assert.NoError(t, err)
	return d
}

func Test_encodedBuffers(t *testing.T) {
	fileBuf, err := newFileBuffer("/tmp/buffer_events.json")
	if err != nil {
		t.Error(err)
	}

	tests := map[string]struct {
		buf encodedBuffer
	}{
		"file buffer": {
			buf: fileBuf,
		},
		"memory buffer": {
			buf: &memBuffer{
				buf: bytes.NewBuffer(make([]byte, 0, 1024)),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := tt.buf
			defer buf.Close()

			writeAndReadBuffer(t, buf, "hello there, 42")
			buf.Reset()

			writeAndReadBuffer(t, buf, "hello there, i'm longer")
			buf.Reset()

			writeAndReadBuffer(t, buf, "then shorter")
			buf.Reset()
		})
	}
}

func writeAndReadBuffer(t *testing.T, buf encodedBuffer, expected string) {
	_, err := buf.Write([]byte(expected))
	if err != nil {
		t.Error(err)
	}

	reader, err := buf.Reader()
	if err != nil {
		t.Error(err)
	}

	result, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, expected, string(result))
}

func ManualTestS3Client(t *testing.T) {
	c, err := NewS3ClientWithConfig(
		S3ClientConfig{
			Config: Config{
				Verbose: true,
			},
			S3: S3{
				Stream:         "tuna",
				Stage:          "dev",
				BufferFilePath: "/tmp/buffer_events.json",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		m := Track{
			Event:  "FooBared",
			UserId: "tuna",
			Properties: map[string]interface{}{
				"index": i,
				"qwer":  3424,
			},
		}
		if err := c.Enqueue(m); err != nil {
			t.Error(err)
		}
	}
	if err := c.Close(); err != nil {
		t.Error(err)
	}

	t.FailNow()
}

func Test_TriggerByTime(t *testing.T) {
	um := uploadMock{
		resultChan: make(chan []byte, 1),
	}
	c := newPatchedClient(t, S3ClientConfig{
		Config: Config{
			Verbose:   true,
			BatchSize: 10,
			Interval:  500 * time.Millisecond,
		},
		S3: S3{
			Stream:         "tuna",
			Stage:          "dev",
			BufferFilePath: "/tmp/buffer_events.tmp",
		},
	}, &um)

	writeEvents(t, c, 1, 0)

	select {
	case msgs := <-um.resultChan:
		msgs1 := readAllMessagies(t, msgs)
		assert.Equal(t, 1, len(msgs1))

	case <-time.After(1 * time.Second):
		t.Errorf("no message by timeout")
	}

	err := c.Close()
	assert.NoError(t, err)

	assert.Empty(t, um.resultChan, "no messages")
}

func Test_MemoryLimit(t *testing.T) {
	const bytesLimit = 5 * 1024 * 1024 // 5 MiB

	um := uploadMock{
		resultChan: make(chan []byte, 1024),
	}
	c := newPatchedClient(t, S3ClientConfig{
		Config: Config{
			Interval:       1 * time.Minute,
			BatchSize:      200_000,
		},
		S3: S3{
			Stream:        "tuna",
			Stage:         "dev",
			MaxBatchBytes: bytesLimit,
			BufferFilePath: "/tmp/buffer_events.tmp",
		},
	}, &um)

	writeEvents(t, c, 400_000, 0)

	readOneEvent(t, um.resultChan, bytesLimit)
	readOneEvent(t, um.resultChan, bytesLimit)

	err := c.Close()
	assert.NoError(t, err)
}

func readOneEvent(t *testing.T, resultChan <-chan []byte, bytesLimit int) {
	eventsData := <-resultChan
	t.Log("event size: ", float64(len(eventsData))/1024/1024, " mib")
	readAllMessagies(t, eventsData)
	assert.GreaterOrEqual(t, len(eventsData), bytesLimit)
}

func Test_MessagesLimit(t *testing.T) {
	const msgsLimit = 10

	um := uploadMock{
		resultChan: make(chan []byte, 1024),
	}
	c := newPatchedClient(t, S3ClientConfig{
		Config: Config{
			Verbose:   true,
			BatchSize: msgsLimit,
		},
		S3: S3{
			Stream:         "tuna",
			Stage:          "dev",
		},
	}, &um)

	writeEvents(t, c, msgsLimit*2, 0)

	msgs1 := readAllMessagies(t, <-um.resultChan)
	assert.Equal(t, msgsLimit, len(msgs1))
	msgs2 := readAllMessagies(t, <-um.resultChan)
	assert.Equal(t, msgsLimit, len(msgs2))

	err := c.Close()
	assert.NoError(t, err)

	assert.Empty(t, um.resultChan, "no messages")
}

func newPatchedClient(t *testing.T, cfg S3ClientConfig, um *uploadMock) Client {
	c, err := NewS3ClientWithConfig(cfg)
	assert.NoError(t, err)

	c.(*s3Client).uploader = um
	return c
}

func readAllMessagies(t *testing.T, input []byte) []Properties {
	gunzip, err := gzip.NewReader(bytes.NewReader(input))
	assert.NoError(t, err)
	defer gunzip.Close()

	scan := bufio.NewReader(gunzip)

	var totalMsgs []Properties
	for {
		row, _, err := scan.ReadLine()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)

		v := Properties{}
		err = json.Unmarshal(row, &v)
		assert.NoError(t, err)

		totalMsgs = append(totalMsgs, v)
	}

	return totalMsgs
}

func writeEvents(t *testing.T, c Client, messages int, delay time.Duration) {
	t.Helper()
	for i := 0; i < messages; i++ {
		m := Track{
			Event:  "FooBared",
			UserId: "tuna",
			Properties: map[string]interface{}{
				"index": i,
				"qwer":  3424,
			},
		}
		if err := c.Enqueue(m); err != nil {
			t.Error(err)
		}
	}
}

type uploadMock struct {
	resultChan chan []byte
}

func (u *uploadMock) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	data, err := ioutil.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}

	u.resultChan <- data

	return nil, nil
}
