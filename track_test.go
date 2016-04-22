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

func TestTrackMissingEvent(t *testing.T) {
	track := Track{
		UserId: "1",
	}

	if err := track.validate(); err == nil {
		t.Error("validating an invalid track object succeeded:", track)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating track:", err)

	} else if e != (FieldError{
		Type:  "analytics.Track",
		Name:  "Event",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating track:", err)
	}
}

func TestTrackMissingUserId(t *testing.T) {
	track := Track{
		Event: "1",
	}

	if err := track.validate(); err == nil {
		t.Error("validating an invalid track object succeeded:", track)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating track:", err)

	} else if e != (FieldError{
		Type:  "analytics.Track",
		Name:  "UserId",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating track:", err)
	}
}

func TestTrackValidWithUserId(t *testing.T) {
	track := Track{
		Event:  "1",
		UserId: "2",
	}

	if err := track.validate(); err != nil {
		t.Error("validating a valid track object failed:", track, err)
	}
}

func TestTrackValidWithAnonymousId(t *testing.T) {
	track := Track{
		Event:       "1",
		AnonymousId: "2",
	}

	if err := track.validate(); err != nil {
		t.Error("validating a valid track object failed:", track, err)
	}
}
