package analytics

var _ Message = GenericMessage(nil)

// This type represents any event type sent in a track call as described in
// https://segment.com/docs/libraries/http/
type GenericMessage map[string]interface{}

func (msg GenericMessage) internal() {
	panic(unimplementedError)
}

func (msg GenericMessage) Validate() error {
	if typ, ok := msg["type"].(string); ok {
		switch typ {
		case "alias":
			m := Alias{Type: "alias"}
			m.UserId, _ = msg["userId"].(string)
			m.PreviousId, _ = msg["previousId"].(string)
			return m.Validate()
		case "group":
			m := Group{Type: "group"}
			m.UserId, _ = msg["userId"].(string)
			m.AnonymousId, _ = msg["anonymousId"].(string)
			m.GroupId, _ = msg["groupId"].(string)
			return m.Validate()
		case "identify":
			m := Identify{Type: "identify"}
			m.UserId, _ = msg["userId"].(string)
			m.AnonymousId, _ = msg["anonymousId"].(string)
			return m.Validate()
		case "page":
			m := Page{Type: "page"}
			m.UserId, _ = msg["userId"].(string)
			m.AnonymousId, _ = msg["anonymousId"].(string)
			return m.Validate()
		case "screen":
			m := Screen{Type: "screen"}
			m.UserId, _ = msg["userId"].(string)
			m.AnonymousId, _ = msg["anonymousId"].(string)
			return m.Validate()
		case "track":
			m := Track{Type: "track"}
			m.UserId, _ = msg["userId"].(string)
			m.AnonymousId, _ = msg["anonymousId"].(string)
			m.Event, _ = msg["event"].(string)
			return m.Validate()
		}
	}
	return FieldError{
		Type:  "analytics.GenericMessage",
		Name:  "Type",
		Value: msg["type"],
	}
}
