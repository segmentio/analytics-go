package analytics

import (
	"reflect"
	"testing"
)

var _ Message = (Event)(nil)

type Event map[string]interface{}

func (e Event) Validate() error {
	return ValidateFields(e)
}

func TestValidateFieldsMissingType(t *testing.T) {
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

func TestValidateFieldsInvalidType(t *testing.T) {
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

func TestValidateFieldsAlias(t *testing.T) {
	msg := Event{
		"type":       "alias",
		"userId":     "user123",
		"previousId": "user456",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic alias message:", err)
	}
}

func TestValidateFieldsAliasInvalid(t *testing.T) {
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

func TestValidateFieldsGroup(t *testing.T) {
	msg := Event{
		"type":    "group",
		"userId":  "user123",
		"groupId": "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestValidateFieldsGroupAnonymous(t *testing.T) {
	msg := Event{
		"type":        "group",
		"anonymousId": "user123",
		"groupId":     "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestValidateFieldsGroupInvalid(t *testing.T) {
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

func TestValidateFieldsIdentify(t *testing.T) {
	msg := Event{
		"type":   "identify",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestValidateFieldsIdentifyAnonymous(t *testing.T) {
	msg := Event{
		"type":        "identify",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestValidateFieldsIdentifyInvalid(t *testing.T) {
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

func TestValidateFieldsPage(t *testing.T) {
	msg := Event{
		"type":   "page",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestValidateFieldsPageAnonymous(t *testing.T) {
	msg := Event{
		"type":        "page",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestValidateFieldsPageInvalid(t *testing.T) {
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

func TestValidateFieldsScreen(t *testing.T) {
	msg := Event{
		"type":   "screen",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestValidateFieldsScreenAnonymous(t *testing.T) {
	msg := Event{
		"type":        "screen",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestValidateFieldsScreenInvalid(t *testing.T) {
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

func TestValidateFieldsTrack(t *testing.T) {
	msg := Event{
		"type":   "track",
		"userId": "user123",
		"event":  "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestValidateFieldsTrackAnonymous(t *testing.T) {
	msg := Event{
		"type":        "track",
		"anonymousId": "user123",
		"event":       "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestValidateFieldsTrackInvalid(t *testing.T) {
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

func TestValidateFieldsQueuePushMaxBatchSize(t *testing.T) {
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

func TestValidateFieldsQueuePushMaxBatchBytes(t *testing.T) {
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
