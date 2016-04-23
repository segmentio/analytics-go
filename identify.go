package analytics

import "time"

// This type represents object sent in an identify call as described in
// https://segment.com/docs/libraries/http/#identify
type Identify struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	Timestamp    time.Time
	Traits       map[string]interface{}
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Identify) validate() error {
	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Identify",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}

func (msg Identify) serializable(msgid string, time time.Time) interface{} {
	return serializableIdentify{
		Type:         "identify",
		MessageId:    makeMessageId(msg.MessageId, msgid),
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		Timestamp:    formatTime(makeTime(msg.Timestamp, time)),
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
	Timestamp    string                 `json:"timestamp,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
