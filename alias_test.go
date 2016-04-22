package analytics

import "testing"

func TestAliasSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	alias := Alias{
		PreviousId: "1",
		UserId:     "2",
		MessageId:  id,
		Timestamp:  ts,
	}

	if v, err := validateSerizable("alias", alias); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
