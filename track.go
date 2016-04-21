package analytics

import "time"

// This type represents object sent in a track call as described in
// https://segment.com/docs/libraries/http/#track
type Track struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	Event        string
	SentAt       time.Time
	Timestamp    time.Time
	Properties   map[string]interface{}
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Track) serializable() interface{} {
	return serializableTrack{
		Type:         "track",
		MessageId:    msg.MessageId,
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		Event:        msg.Event,
		SentAt:       formatTime(msg.SentAt),
		Timestamp:    formatTime(msg.Timestamp),
		Context:      msg.Context,
		Integrations: msg.Integrations,
		Properties:   msg.Properties,
	}
}

type serializableTrack struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Event        string                 `json:"event"`
	SentAt       string                 `json:"sentAt,omitempty"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
