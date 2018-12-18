package analytics

import (
	"encoding/json"
	"strconv"
	"time"
)

// Time is a wrapper over time.Time with custom json marshalling.
type Time time.Time

// MarshalJSON marshals time as unix milliseconds.
func (t Time) MarshalJSON() ([]byte, error) {
	if t == Time(time.Time{}) {
		return json.Marshal(0)
	}
	millis := int64(time.Millisecond / time.Nanosecond)
	return json.Marshal(time.Time(t).UnixNano() / millis)
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
