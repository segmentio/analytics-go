package analytics

import "testing"

func TestMessageIdDefault(t *testing.T) {
	if id := makeMessageId("", "42"); id != "42" {
		t.Error("invalid default message id:", id)
	}
}

func TestMessageIdNonDefault(t *testing.T) {
	if id := makeMessageId("A", "42"); id != "A" {
		t.Error("invalid non-default message id:", id)
	}
}
