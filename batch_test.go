package analytics

import "testing"

func TestBatchSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	batch := batch{
		MessageId: id,
		SentAt:    ts,
		Messages: []interface{}{
			(Track{
				UserId: "1",
				Event:  "wake-up",
			}).serializable(id, ts),

			(Group{
				UserId:  "1",
				GroupId: "2",
			}).serializable(id, ts),

			(Page{
				UserId: "1",
				Name:   "home",
			}).serializable(id, ts),

			(Alias{
				PreviousId: "1",
				UserId:     "2",
			}).serializable(id, ts),

			(Identify{
				UserId: "1",
			}).serializable(id, ts),
		},
	}

	if v, err := validateSerizable("", batch); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}
