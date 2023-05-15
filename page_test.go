package journify

import "testing"

func TestPageMissingUserId(t *testing.T) {
	page := Page{}

	if err := page.Validate(); err == nil {
		t.Error("validating an invalid page object succeeded:", page)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating page:", err)

	} else if e != (FieldError{
		Type:  "journify.Page",
		Name:  "UserId",
		Value: "",
	}) {
		t.Error("invalid error value returned when validating page:", err)
	}
}

func TestPageValidWithUserId(t *testing.T) {
	page := Page{
		UserId: "2",
	}

	if err := page.Validate(); err != nil {
		t.Error("validating a valid page object failed:", page, err)
	}
}

func TestPageValidWithAnonymousId(t *testing.T) {
	page := Page{
		AnonymousId: "2",
	}

	if err := page.Validate(); err != nil {
		t.Error("validating a valid page object failed:", page, err)
	}
}
