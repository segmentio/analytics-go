package analytics

import (
	"encoding/json"
	"time"
)

// Values implementing this interface are used by analytics clients to notify
// the application when a message send succeeded or failed.
//
// Callback methods are called by a client's internal goroutines, there are no
// guarantees on which goroutine will trigger the callbacks, the calls can be
// made sequentially or in parallel, the order doesn't depend on the order of
// messages were queued to the client.
//
// Callback methods must return quickly and not cause long blocking operations
// to avoid interferring with the client's internal work flow.
type Callback interface {

	// This method is called for every message that was successfully sent to
	// the API.
	Success(Message)

	// This method is called for every message that failed to be sent to the
	// API and will be discarded by the client.
	Failure(Message, error)
}

const (
	messageIDKey = "messageId"
	timestampKey = "timestamp"
)

type Message map[string]interface{}

// EnsureMessageID sets the message ID if it hasn't already been set.
func (m Message) EnsureMessageID(id string) {
	if _, ok := m[messageIDKey]; !ok {
		m[messageIDKey] = id
	}
}

// EnsureTimestamp sets the timestamp if it hasn't already been set.
func (m Message) EnsureTimestamp(ts time.Time) {
	if _, ok := m[timestampKey]; !ok {
		m[timestampKey] = ts
	}
}

// This structure represents objects sent to the /v1/batch endpoint. We don't
// export this type because it's only meant to be used internally to send groups
// of messages in one API call.
type batch struct {
	MessageId string                 `json:"messageId"`
	SentAt    time.Time              `json:"sentAt"`
	Messages  []message              `json:"batch"`
	Context   map[string]interface{} `json:"context"`
}

type message struct {
	msg  Message
	json []byte
}

func makeMessage(m Message, maxBytes int) (msg message, err error) {
	if msg.json, err = json.Marshal(m); err == nil {
		if len(msg.json) > maxBytes {
			err = ErrMessageTooBig
		} else {
			msg.msg = m
		}
	}
	return
}

func (m message) MarshalJSON() ([]byte, error) {
	return m.json, nil
}

func (m message) size() int {
	// The `+ 1` is for the comma that sits between each items of a JSON array.
	return len(m.json) + 1
}

type messageQueue struct {
	pending       []message
	bytes         int
	maxBatchSize  int
	maxBatchBytes int
}

func (q *messageQueue) push(m message) (b []message) {
	if (q.bytes + m.size()) > q.maxBatchBytes {
		b = q.flush()
	}

	if q.pending == nil {
		q.pending = make([]message, 0, q.maxBatchSize)
	}

	q.pending = append(q.pending, m)
	q.bytes += len(m.json)

	if b == nil && len(q.pending) == q.maxBatchSize {
		b = q.flush()
	}

	return
}

func (q *messageQueue) flush() (msgs []message) {
	msgs, q.pending, q.bytes = q.pending, nil, 0
	return
}

const (
	maxBatchBytes   = 500000
	maxMessageBytes = 32000
)
