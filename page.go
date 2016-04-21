package analytics

import "time"

// This type represents object sent in a page call as described in
// https://segment.com/docs/libraries/http/#page
type Page struct {
	MessageId    string
	AnonymousId  string `json:"anonymousId,omitempty"`
	UserId       string `json:"userId,omitempty"`
	Name         string `json:"name,omitempty"`
	SentAt       time.Time
	Timestamp    time.Time
	Traits       map[string]interface{} `json:"properties,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}

func (msg Page) serializable() interface{} {
	return serializablePage{
		Type:         "page",
		MessageId:    msg.MessageId,
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		Name:         msg.Name,
		SentAt:       formatTime(msg.SentAt),
		Timestamp:    formatTime(msg.Timestamp),
		Traits:       msg.Traits,
		Context:      msg.Context,
		Integrations: msg.Integrations,
	}
}

type serializablePage struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	Name         string                 `json:"name,omitempty"`
	SentAt       string                 `json:"sentAt,omitempty"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Traits       map[string]interface{} `json:"properties,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
