package journify

type FieldGetter interface {
	GetField(field string) (interface{}, bool)
}

func getString(msg FieldGetter, field string) string {
	if val, ok := msg.GetField(field); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func ValidateFields(msg FieldGetter) error {
	typ, _ := msg.GetField("type")
	if str, ok := typ.(string); ok {
		switch str {
		case "group":
			return Group{
				Type:        "group",
				UserId:      getString(msg, "userId"),
				AnonymousId: getString(msg, "anonymousId"),
				GroupId:     getString(msg, "groupId"),
			}.Validate()
		case "identify":
			return Identify{
				Type:        "identify",
				UserId:      getString(msg, "userId"),
				AnonymousId: getString(msg, "anonymousId"),
			}.Validate()
		case "page":
			return Page{
				Type:        "page",
				UserId:      getString(msg, "userId"),
				AnonymousId: getString(msg, "anonymousId"),
			}.Validate()
		case "track":
			return Track{
				Type:        "track",
				UserId:      getString(msg, "userId"),
				AnonymousId: getString(msg, "anonymousId"),
				Event:       getString(msg, "event"),
			}.Validate()
		}
	}
	return FieldError{
		Type:  "journify.Event",
		Name:  "Type",
		Value: typ,
	}
}
