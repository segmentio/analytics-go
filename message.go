package analytics

import "time"

// This interface is used to represent analytics objects that can be sent via
// a client.
//
// Types like analytics.Track, analytics.Page, etc... implement this interface
// and therefore can be passed to the analytics.Client.Send method.
type Message interface {

	// Validates the internal structure of the message, the method must return
	// nil if the message is valid, or an error describing what went wrong.
	validate() error
}

// Takes a message id as first argument and returns it, unless it's the zero-
// value, in that case the default id passed as second argument is returned.
func makeMessageId(id string, def string) string {
	if len(id) == 0 {
		return def
	}
	return id
}

// Returns a string representation of the time value passed as first argument,
// unless it's a zero-value, in that case the default value passed as second
// argument is used instead.
func makeTimestamp(t time.Time, def time.Time) time.Time {
	if t == (time.Time{}) {
		return def
	}
	return t
}

// This structure represents objects sent to the /v1/batch endpoint. We don't
// export this type because it's only meant to be used internally to send groups
// of messages in one API call.
type batch struct {
	MessageId string    `json:"messageId"`
	SentAt    time.Time `json:"sentAt"`
	Messages  []Message `json:"batch"`
	Context   *Context  `json:"context"`
}
