package analytics

import "time"

const (
	iso8601 = "2006-01-02T15:04:05.999-0700"
)

// Formats a time value to an iso8601 representation is millisecond precision.
//
// We use this function instead of relying on the default serialization of
// json.Marshal because since time.Time is a struct type the `omitempty` tag is
// ignored and the zero value ends up being serialized which causes the message
// to have an invalid timestamp of January 1st 1970.
//
// This also allows us to control the format in which the time is formatted and
// ensure it is the same that has been historically used.
func formatTime(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}
	return t.Format(iso8601)
}

// Returns the time value passed as first argument, unless it's the zero-value,
// in that case the default value passed as second argument is returned.
func makeTime(t time.Time, def time.Time) time.Time {
	if t == (time.Time{}) {
		return def
	}
	return t
}
