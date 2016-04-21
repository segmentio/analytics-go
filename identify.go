package analytics

import "time"

// This type represents object sent in an identify call as described in
// https://segment.com/docs/libraries/http/#identify
type Identify struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	SentAt       time.Time
	Timestamp    time.Time
	Traits       map[string]interface{}
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Identify) serializable() interface{} {
	return serializableIdentify{
		Type:         "identify",
		MessageId:    msg.MessageId,
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		SentAt:       formatTime(self.SentAt),
		Timestamp:    formatTime(self.Timestamp),
		Traits:       msg.Traits,
		Context:      msg.Context,
		Integrations: msg.Integrations,
	}
}

type serializableIdentify struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	SentAt       string                 `json:"sentAt,omitempty"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
