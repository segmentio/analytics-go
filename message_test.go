package analytics

import (
	"reflect"
	"testing"
)

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

func TestMessageQueuePushMaxBatchSize(t *testing.T) {
	m0, _ := makeMessage(Track{
		UserId: "1",
		Event:  "A",
	})

	m1, _ := makeMessage(Track{
		UserId: "2",
		Event:  "A",
	})

	q := messageQueue{
		maxBatchSize:  2,
		maxBatchBytes: maxBatchBytes,
	}

	if msgs := q.push(m0); msgs != nil {
		t.Error("unexpected message batch returned after pushing only one message")
	}

	if msgs := q.push(m1); !reflect.DeepEqual(msgs, []message{m0, m1}) {
		t.Error("invalid message batch returned after pushing two messages:", msgs)
	}
}

func TestMessageQueuePushMaxBatchBytes(t *testing.T) {
	m0, _ := makeMessage(Track{
		UserId: "1",
		Event:  "A",
	})

	m1, _ := makeMessage(Track{
		UserId: "2",
		Event:  "A",
	})

	q := messageQueue{
		maxBatchSize:  100,
		maxBatchBytes: len(m0.json) + 1,
	}

	if msgs := q.push(m0); msgs != nil {
		t.Error("unexpected message batch returned after pushing only one message")
	}

	if msgs := q.push(m1); !reflect.DeepEqual(msgs, []message{m0}) {
		t.Error("invalid message batch returned after pushing two messages:", msgs)
	}

	if !reflect.DeepEqual(q.pending, []message{m1}) {
		t.Error("invalid state of the message queue after pushing two messages:", q.pending)
	}
}
