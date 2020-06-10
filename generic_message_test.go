package analytics

import "testing"

func TestGenericMessageMissingType(t *testing.T) {
	msg := GenericMessage{
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.GenericMessage",
		Name:  "Type",
		Value: nil,
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageInvalidType(t *testing.T) {
	msg := GenericMessage{
		"type":   "invalid",
		"userId": "user123",
	}

	if err := msg.Validate(); err == nil {
		t.Error("validating an invalid generic message succeeded:", msg)
	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating a generic message:", err)

	} else if e != (FieldError{
		Type:  "analytics.GenericMessage",
		Name:  "Type",
		Value: "invalid",
	}) {
		t.Error("invalid error type returned when validating a generic message:", err)
	}
}

func TestGenericMessageAlias(t *testing.T) {
	msg := GenericMessage{
		"type":       "alias",
		"userId":     "user123",
		"previousId": "user456",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic alias message:", err)
	}
}

func TestGenericMessageAliasInvalid(t *testing.T) {
	msg := GenericMessage{
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
	msg := GenericMessage{
		"type":    "group",
		"userId":  "user123",
		"groupId": "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestGenericMessageGroupAnonymous(t *testing.T) {
	msg := GenericMessage{
		"type":        "group",
		"anonymousId": "user123",
		"groupId":     "group1",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic group message:", err)
	}
}

func TestGenericMessageGroupInvalid(t *testing.T) {
	msg := GenericMessage{
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
	msg := GenericMessage{
		"type":   "identify",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestGenericMessageIdentifyAnonymous(t *testing.T) {
	msg := GenericMessage{
		"type":        "identify",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic identify message:", err)
	}
}

func TestGenericMessageIdentifyInvalid(t *testing.T) {
	msg := GenericMessage{
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
	msg := GenericMessage{
		"type":   "page",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestGenericMessagePageAnonymous(t *testing.T) {
	msg := GenericMessage{
		"type":        "page",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic page message:", err)
	}
}

func TestGenericMessagePageInvalid(t *testing.T) {
	msg := GenericMessage{
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
	msg := GenericMessage{
		"type":   "screen",
		"userId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestGenericMessageScreenAnonymous(t *testing.T) {
	msg := GenericMessage{
		"type":        "screen",
		"anonymousId": "user123",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic screen message:", err)
	}
}

func TestGenericMessageScreenInvalid(t *testing.T) {
	msg := GenericMessage{
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
	msg := GenericMessage{
		"type":   "track",
		"userId": "user123",
		"event":  "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestGenericMessageTrackAnonymous(t *testing.T) {
	msg := GenericMessage{
		"type":        "track",
		"anonymousId": "user123",
		"event":       "testing",
	}
	if err := msg.Validate(); err != nil {
		t.Error("error returned when validating a generic track message:", err)
	}
}

func TestGenericMessageTrackInvalid(t *testing.T) {
	msg := GenericMessage{
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
