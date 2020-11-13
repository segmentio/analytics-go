package analytics

// Track represents object sent in a track call as described in
// https://segment.com/docs/libraries/http/#track
type Track struct {
	// This field is exported for serialization purposes and shouldn't be set by
	// the application, its value is always overwritten by the library.
	Type string `json:"type,omitempty"`

	MessageId    string       `json:"messageId,omitempty"`
	AnonymousId  string       `json:"anonymousId,omitempty"`
	UserId       string       `json:"userId,omitempty"`
	Event        string       `json:"event"`
	Timestamp    Time         `json:"timestamp,omitempty"`
	Context      *Context     `json:"context,omitempty"`
	Properties   Properties   `json:"properties,omitempty"`
	Integrations Integrations `json:"integrations,omitempty"`
}

func (msg Track) tags() []string {
	return []string{"type:" + msg.Type, "event:" + msg.Event}
}

func (msg Track) validate() error {
	if len(msg.Event) == 0 {
		return FieldError{
			Type:  "analytics.Track",
			Name:  "Event",
			Value: msg.Event,
		}
	}

	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Track",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}

// TrackObj represents object sent in a track call as Track
// but instead of map[string]interface{} accepts any struct which should be serialized to json
type TrackObj struct {
	// This field is exported for serialization purposes and shouldn't be set by
	// the application, its value is always overwritten by the library.
	Track
	Properties interface{} `json:"properties,omitempty"`
}

// TrackObjLess represents object sent in a track call as TrackObj
// but only an event name is mandatory field
type TrackObjLess struct {
	// This field is exported for serialization purposes and shouldn't be set by
	// the application, its value is always overwritten by the library.
	Track
	Properties interface{} `json:"properties,omitempty"`
}

func (msg TrackObjLess) validate() error {
	if len(msg.Event) == 0 {
		return FieldError{
			Type:  "analytics.Track",
			Name:  "Event",
			Value: msg.Event,
		}
	}

	return nil
}
