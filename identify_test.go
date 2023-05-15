package journify

import "testing"

func TestIdentifyMissingUserId(t *testing.T) {
	identify := Identify{}

	if err := identify.Validate(); err == nil {
		t.Error("validating an invalid identify object succeeded:", identify)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating identify:", err)

	} else if e != (FieldError{
		Type:  "journify.Identify",
		Name:  "UserId",
		Value: "",
	}) {
		t.Error("invalid error value returned when validating identify:", err)
	}
}

func TestIdentifyValidWithUserId(t *testing.T) {
	identify := Identify{
		UserId: "2",
	}

	if err := identify.Validate(); err != nil {
		t.Error("validating a valid identify object failed:", identify, err)
	}
}

func TestIdentifyValidWithAnonymousId(t *testing.T) {
	identify := Identify{
		AnonymousId: "2",
	}

	if err := identify.Validate(); err != nil {
		t.Error("validating a valid identify object failed:", identify, err)
	}
}
