package analytics

import "testing"

func TestParseJsonTagEmpty(t *testing.T) {
	name, omitempty := parseJsonTag("", "default")

	if name != "default" {
		t.Error("invalid field name found in empty tag:", name)
	}

	if omitempty {
		t.Error("unexpected 'omitempty' state found in empty tag")
	}
}

func TestParseJsonTagName(t *testing.T) {
	name, omitempty := parseJsonTag("name", "default")

	if name != "name" {
		t.Error("invalid field name found in json tag:", name)
	}

	if omitempty {
		t.Error("unexpected 'omitempty' state found in json tag")
	}
}

func TestParseJsonTagOmitempty(t *testing.T) {
	name, omitempty := parseJsonTag(",omitempty", "default")

	if name != "default" {
		t.Error("invalid field name found in omitempty tag:", name)
	}

	if !omitempty {
		t.Error("expected 'omitempty' state not found in json tag")
	}
}

func TestParseJsonTagNameOmitempty(t *testing.T) {
	name, omitempty := parseJsonTag("name,omitempty", "default")

	if name != "name" {
		t.Error("invalid field name found in json tag:", name)
	}

	if !omitempty {
		t.Error("expected 'omitempty' state not found in json tag")
	}
}
