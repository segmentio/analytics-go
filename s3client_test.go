package analytics

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
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

	encoder := &bufferedEncoder{
		maxBatchSize:  100,
		maxBatchBytes: 10,
		newBufFunc: func() encodedBuffer {
			return newMemBuffer(1024)
		},
	}
	encoder.init()

	_, err := encodeMessage(encoder, m, nil, func() Time { return Time{} })
	require.NoError(t, err)

	buf, err := encoder.CommitBuffer()
	require.NoError(t, err)
	reader, err := buf.Reader()
	require.NoError(t, err)

	result := readAndUngzip(t, reader)

	expected := `{"event":{"userId":"tuna","event":"FooBared","timestamp":0,"properties":{"index":1,"qwer":3424}},"sentAt":0,"receivedAt":0}` + "\n"

	require.Equal(t, expected, string(result))
}

func readAndUngzip(t *testing.T, r io.Reader) []byte {
	q, err := gzip.NewReader(r)
	require.NoError(t, err)
	defer q.Close()

	d, err := ioutil.ReadAll(q)
	require.NoError(t, err)
	return d
}

func Test_encodedBuffers(t *testing.T) {
	fileBuf, err := newFileBuffer(filePath)
	require.NoError(t, err)
	defer fileBuf.Close()

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
	require.NoError(t, err)

	reader, err := buf.Reader()
	require.NoError(t, err)

	result, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	require.Equal(t, expected, string(result))
}

const (
	filePath = "/tmp/buffer_events.tmp"
)

func ManualTestS3Client(t *testing.T) {
	c, err := NewS3ClientWithConfig(
		S3ClientConfig{
			Config: Config{
				Verbose: true,
			},
			S3: S3{
				Stream:         "tuna",
				Stage:          "dev",
				BufferFilePath: filePath,
			},
		},
	)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		m := Track{
			Event:  "FooBared",
			UserId: "tuna",
			Properties: map[string]interface{}{
				"index": i,
				"qwer":  3424,
			},
		}
		require.NoError(t, c.Enqueue(m))
	}
	require.NoError(t, c.Close())

	t.FailNow()
}

func Test_TriggerByTime(t *testing.T) {
	defer checkNoFilesLeft(t, filePath)
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
			BufferFilePath: filePath,
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

	require.NoError(t, c.Close())
	require.Empty(t, um.resultChan, "no messages")
}

func Test_MemoryLimit(t *testing.T) {
	defer checkNoFilesLeft(t, filePath)
	const bytesLimit = 5 * 1024 * 1024 // 5 MiB

	um := uploadMock{
		resultChan: make(chan []byte, 4),
	}
	c := newPatchedClient(t, S3ClientConfig{
		Config: Config{
			Interval:  5 * time.Minute,
			BatchSize: 200_000,
		},
		S3: S3{
			Stream:         "tuna",
			Stage:          "dev",
			MaxBatchBytes:  bytesLimit,
			BufferFilePath: filePath,
		},
	}, &um)

	writeEvents(t, c, 400_000, 0)

	readOneEvent(t, um.resultChan, bytesLimit)
	readOneEvent(t, um.resultChan, bytesLimit)

	require.NoError(t, c.Close())
}

func checkNoFilesLeft(t *testing.T, path string) {
	t.Helper()
	dir, fn := filepath.Split(path)
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		assert.NotContains(t, file.Name(), fn)
	}
}

func readOneEvent(t *testing.T, resultChan <-chan []byte, bytesLimit int) {
	eventsData := <-resultChan
	log.Println("event size: ", float64(len(eventsData))/1024/1024, " mib")
	readAllMessagies(t, eventsData)
	require.GreaterOrEqual(t, len(eventsData), bytesLimit)
}

func Test_MessagesLimit(t *testing.T) {
	const msgsLimit = 10

	um := uploadMock{
		resultChan: make(chan []byte, 3),
	}
	c := newPatchedClient(t, S3ClientConfig{
		Config: Config{
			Verbose:               true,
			BatchSize:             msgsLimit,
			maxConcurrentRequests: 1,
		},
		S3: S3{
			Stream: "tuna",
			Stage:  "dev",
		},
	}, &um)

	writeEvents(t, c, msgsLimit*2, 0)

	msgs1 := readAllMessagies(t, <-um.resultChan)
	require.Equal(t, msgsLimit, len(msgs1))
	msgs2 := readAllMessagies(t, <-um.resultChan)
	require.Equal(t, msgsLimit, len(msgs2))

	err := c.Close()
	require.NoError(t, err)

	require.Empty(t, um.resultChan, "no messages")
}

func newPatchedClient(t *testing.T, cfg S3ClientConfig, um *uploadMock) Client {
	c, err := NewS3ClientWithConfig(cfg)
	require.NoError(t, err)

	c.(*s3Client).uploader = um
	return c
}

func readAllMessagies(t *testing.T, input []byte) []Properties {
	gunzip, err := gzip.NewReader(bytes.NewReader(input))
	require.NoError(t, err)
	defer gunzip.Close()

	scan := bufio.NewReader(gunzip)

	var totalMsgs []Properties
	for {
		row, _, err := scan.ReadLine()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		v := Properties{}
		err = json.Unmarshal(row, &v)
		require.NoError(t, err)

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
		require.NoError(t, c.Enqueue(m))
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
