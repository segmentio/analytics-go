package analytics

func getString(msg map[string]interface{}, field string) string {
	val, _ := msg[field].(string)
	return val
}

func ValidateFields(msg map[string]interface{}) error {
	if typ, ok := msg["type"].(string); ok {
		switch typ {
		case "alias":
			return Alias{
				Type:       "alias",
				UserId:     getString(msg, "userId"),
				PreviousId: getString(msg, "previousId"),
			}.Validate()
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
		case "screen":
			return Screen{
				Type:        "screen",
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
		Type:  "analytics.Event",
		Name:  "Type",
		Value: msg["type"],
	}
}
