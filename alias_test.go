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

func TestAliasMissingUserId(t *testing.T) {
	alias := Alias{
		PreviousId: "1",
	}

	if err := alias.validate(); err == nil {
		t.Error("validating an invalid alias object succeeded:", alias)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating alias:", err)

	} else if e != (FieldError{
		Type:  "analytics.Alias",
		Name:  "UserId",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating alias:", err)
	}
}

func TestAliasMissingPreviousId(t *testing.T) {
	alias := Alias{
		UserId: "1",
	}

	if err := alias.validate(); err == nil {
		t.Error("validating an invalid alias object succeeded:", alias)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating alias:", err)

	} else if e != (FieldError{
		Type:  "analytics.Alias",
		Name:  "PreviousId",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating alias:", err)
	}
}

func TestAliasValid(t *testing.T) {
	alias := Alias{
		PreviousId: "1",
		UserId:     "2",
	}

	if err := alias.validate(); err != nil {
		t.Error("validating a valid alias object failed:", alias, err)
	}
}
