package analytics

import (
	"strconv"
	"time"
)

// Time is a wrapper over time.Time with custom json marshalling.
type Time time.Time

// MarshalJSON marshals time as unix milliseconds.
func (t Time) MarshalJSON() ([]byte, error) {
	var ts int64
	if t != Time(time.Time{}) {
		ts = time.Time(t).UnixNano() / int64(time.Millisecond)
	}

	return strconv.AppendInt(nil, ts, 10), nil
}

// UnmarshalJSON unmarshals milliseconds into Time.
func (t *Time) UnmarshalJSON(data []byte) error {
	millis, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(0, millis*int64(time.Millisecond)))
	return nil
}
