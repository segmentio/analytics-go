package analytics

import "testing"

func TestIdentifySerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	identify := Identify{
		UserId:    "1",
		MessageId: id,
		Timestamp: ts,
	}

	if v, err := validateSerizable("identify", identify); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
