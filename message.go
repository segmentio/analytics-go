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

	// Returns a serializable representation of the message, using the given
	// message id and timestamp pass as argument if none were already set on
	// the message.
	serializable(msgid string, time time.Time) interface{}
}

// Takes a message id as first argument and returns it, unless it's the zero-
// value, in that case the default id passed as second argument is returned.
func makeMessageId(id string, def string) string {
	if len(id) == 0 {
		return def
	}
	return id
}
