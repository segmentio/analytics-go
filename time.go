package analytics

import "time"

const (
	iso8601 = "2006-01-02T15:04:05.999Z07:00"
)

// Formats a time value to an iso8601 representation is millisecond precision.
// We use this function instead of relying on the default serialization of
// json.Marshal because since time.Time is a struct type the `omitempty` tag is
// ignored and the zero value ends up being serialized which causes the message
// to have an invalid timestamp of January 1st 1970.
func formatTime(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}
	return t.Format(iso8601)
}
