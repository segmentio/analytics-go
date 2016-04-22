package analytics

import "testing"

func TestGroupSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	group := Group{
		UserId:    "1",
		GroupId:   "2",
		MessageId: id,
		Timestamp: ts,
	}

	if v, err := validateSerizable("group", group); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
