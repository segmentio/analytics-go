package analytics

import "testing"

func TestPageSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	page := Page{
		UserId:    "1",
		Name:      "home",
		MessageId: id,
		Timestamp: ts,
	}

	if v, err := validateSerizable("page", page); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
