package analytics

import "time"

// This type represents object sent in a page call as described in
// https://segment.com/docs/libraries/http/#page
type Page struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	Name         string
	Timestamp    time.Time
	Context      Context
	Properties   map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Page) validate() error {
	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Page",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}

func (msg Page) serializable(msgid string, time time.Time) interface{} {
	return serializablePage{
		Type:         "page",
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

type serializablePage struct {
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
