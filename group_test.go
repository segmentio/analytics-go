package analytics

import "testing"

func TestGroupSerializable(t *testing.T) {
	id := mockId()
	ts := mockTime()

	group := Group{
		UserId:    "1",
		GroupId:   "2",
		MessageId: id,
		Timestamp: ts,
	}

	if v, err := validateSerizable("group", group); err != nil {
		t.Errorf("%s: %#v", err, v)
	}
}

func TestGroupMissingGroupId(t *testing.T) {
	group := Group{
		UserId: "1",
	}

	if err := group.validate(); err == nil {
		t.Error("validating an invalid group object succeeded:", group)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating group:", err)

	} else if e != (FieldError{
		Type:  "analytics.Group",
		Name:  "GroupId",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating group:", err)
	}
}

func TestGroupMissingUserId(t *testing.T) {
	group := Group{
		GroupId: "1",
	}

	if err := group.validate(); err == nil {
		t.Error("validating an invalid group object succeeded:", group)

	} else if e, ok := err.(FieldError); !ok {
		t.Error("invalid error type returned when validating group:", err)

	} else if e != (FieldError{
		Type:  "analytics.Group",
		Name:  "UserId",
		Value: "",
	}) {
		t.Errorf("invalid error value returned when validating group:", err)
	}
}

func TestGroupValidWithUserId(t *testing.T) {
	group := Group{
		GroupId: "1",
		UserId:  "2",
	}

	if err := group.validate(); err != nil {
		t.Error("validating a valid group object failed:", group, err)
	}
}

func TestGroupValidWithAnonymousId(t *testing.T) {
	group := Group{
		GroupId:     "1",
		AnonymousId: "2",
	}

	if err := group.validate(); err != nil {
		t.Error("validating a valid group object failed:", group, err)
	}
}
