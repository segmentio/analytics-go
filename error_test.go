package journify

import "testing"

func TestConfigError(t *testing.T) {
	e := ConfigError{
		Reason: "testing",
		Field:  "Answer",
		Value:  42,
	}

	if s := e.Error(); s != "journify.NewWithConfig: testing (journify.Config.Answer: 42)" {
		t.Error("invalid error message returned by config error:", s)
	}
}

func TestFieldError(t *testing.T) {
	e := FieldError{
		Type:  "testing.T",
		Name:  "Answer",
		Value: 42,
	}

	if s := e.Error(); s != "testing.T.Answer: invalid field value: 42" {
		t.Error("invalid error message returned by field error:", s)
	}
}
