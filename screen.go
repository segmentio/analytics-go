package analytics

import "time"

// This type represents object sent in a screen call as described in
// https://segment.com/docs/libraries/http/#screen
type Screen struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	Name         string
	Timestamp    time.Time
	Context      Context
	Properties   map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Screen) validate() error {
	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Screen",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}

func (msg Screen) serializable(msgid string, time time.Time) interface{} {
	return serializableScreen{
		Type:         "screen",
		MessageId:    makeMessageId(msg.MessageId, msgid),
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		Name:         msg.Name,
		Timestamp:    makeTimestamp(msg.Timestamp, time),
		Context:      makeJsonContext(msg.Context),
		Properties:   msg.Properties,
		Integrations: msg.Integrations,
	}
}

type serializableScreen struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Context      *Context               `json:"context,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
