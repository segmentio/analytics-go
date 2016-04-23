package analytics

import "testing"

func TestParseJsonTagEmpty(t *testing.T) {
	name, omitempty := parseJsonTag("")

	if name != "" {
		t.Error("invalid field name found in empty tag:", name)
	}

	if omitempty {
		t.Error("unexpected 'omitempty' state found in empty tag")
	}
}

func TestParseJsonTagName(t *testing.T) {
	name, omitempty := parseJsonTag("name")

	if name != "name" {
		t.Error("invalid field name found in json tag:", name)
	}

	if omitempty {
		t.Error("unexpected 'omitempty' state found in json tag")
	}
}

func TestParseJsonTagOmitempty(t *testing.T) {
	name, omitempty := parseJsonTag(",omitempty")

	if name != "" {
		t.Error("invalid field name found in omitempty tag:", name)
	}

	if !omitempty {
		t.Error("expected 'omitempty' state not found in json tag")
	}
}

func TestParseJsonTagNameOmitempty(t *testing.T) {
	name, omitempty := parseJsonTag("name,omitempty")

	if name != "name" {
		t.Error("invalid field name found in json tag:", name)
	}

	if !omitempty {
		t.Error("expected 'omitempty' state not found in json tag")
	}
}
