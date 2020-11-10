package analytics

import (
	"reflect"
	"testing"
)

func TestGenericMessageMissingType(t *testing.T) {
	msg := Event{
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Event",
		Name:  "Type",
		Value: nil,
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageInvalidType(t *testing.T) {
	msg := Event{
		"type":   "invalid",
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Event",
		Name:  "Type",
		Value: "invalid",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageAlias(t *testing.T) {
	msg := Event{
		"type":       "alias",
		"userId":     "user123",
		"previousId": "user456",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic alias message:", err)
	}
}

func TestGenericMessageAliasInvalid(t *testing.T) {
	msg := Event{
		"type":   "alias",
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Alias",
		Name:  "PreviousId",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageGroup(t *testing.T) {
	msg := Event{
		"type":    "group",
		"userId":  "user123",
		"groupId": "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestGenericMessageGroupAnonymous(t *testing.T) {
	msg := Event{
		"type":        "group",
		"anonymousId": "user123",
		"groupId":     "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestGenericMessageGroupInvalid(t *testing.T) {
	msg := Event{
		"type":   "group",
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Group",
		Name:  "GroupId",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageIdentify(t *testing.T) {
	msg := Event{
		"type":   "identify",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestGenericMessageIdentifyAnonymous(t *testing.T) {
	msg := Event{
		"type":        "identify",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestGenericMessageIdentifyInvalid(t *testing.T) {
	msg := Event{
		"type": "identify",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Identify",
		Name:  "UserId",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessagePage(t *testing.T) {
	msg := Event{
		"type":   "page",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestGenericMessagePageAnonymous(t *testing.T) {
	msg := Event{
		"type":        "page",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestGenericMessagePageInvalid(t *testing.T) {
	msg := Event{
		"type": "page",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Page",
		Name:  "UserId",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageScreen(t *testing.T) {
	msg := Event{
		"type":   "screen",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestGenericMessageScreenAnonymous(t *testing.T) {
	msg := Event{
		"type":        "screen",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestGenericMessageScreenInvalid(t *testing.T) {
	msg := Event{
		"type": "screen",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Screen",
		Name:  "UserId",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageTrack(t *testing.T) {
	msg := Event{
		"type":   "track",
		"userId": "user123",
		"event":  "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestGenericMessageTrackAnonymous(t *testing.T) {
	msg := Event{
		"type":        "track",
		"anonymousId": "user123",
		"event":       "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestGenericMessageTrackInvalid(t *testing.T) {
	msg := Event{
		"type":   "track",
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.Track",
		Name:  "Event",
		Value: "",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageQueuePushMaxBatchSize(t *testing.T) {
	m0, _ := makeMessage(Event{
		"type":   "track",
		"userId": "1",
		"event":  "A",
	}, maxMessageBytes)

	m1, _ := makeMessage(Event{
		"type":   "track",
		"userId": "2",
		"event":  "A",
	}, maxMessageBytes)

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

func TestGenericMessageQueuePushMaxBatchBytes(t *testing.T) {
	m0, _ := makeMessage(Event{
		"type":   "track",
		"UserId": "1",
		"Event":  "A",
	}, maxMessageBytes)

	m1, _ := makeMessage(Event{
		"type":   "track",
		"UserId": "2",
		"Event":  "A",
	}, maxMessageBytes)

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
