package analytics

import "time"

// This type represents object sent in a track call as described in
// https://segment.com/docs/libraries/http/#track
type Track struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	Event        string
	Timestamp    time.Time
	Context      Context
	Properties   map[string]interface{}
	Integrations map[string]interface{}
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

func (msg Track) serializable(msgid string, time time.Time) interface{} {
	return serializableTrack{
		Type:         "track",
		MessageId:    makeMessageId(msg.MessageId, msgid),
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		Event:        msg.Event,
		Timestamp:    makeTimestamp(msg.Timestamp, time),
		Context:      makeJsonContext(msg.Context),
		Properties:   msg.Properties,
		Integrations: msg.Integrations,
	}
}

type serializableTrack struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Event        string                 `json:"event"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Context      *Context               `json:"context,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
