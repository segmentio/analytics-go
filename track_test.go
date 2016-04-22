package analytics

import "testing"

func TestTrackSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	track := Track{
		UserId:    "1",
		Event:     "wake-up",
		MessageId: id,
		Timestamp: ts,
	}

	if v, err := validateSerizable("track", track); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
