package analytics

import "time"

// This type represents object sent in a alias call as described in
// https://segment.com/docs/libraries/http/#alias
type Alias struct {
	MessageId    string
	PreviousId   string
	UserId       string
	Timestamp    time.Time
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Alias) validate() error {
	if len(msg.UserId) == 0 {
		return FieldError{
			Type:  "analytics.Alias",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	if len(msg.PreviousId) == 0 {
		return FieldError{
			Type:  "analytics.Alias",
			Name:  "PreviousId",
			Value: msg.PreviousId,
		}
	}

	return nil
}

func (msg Alias) serializable(msgid string, time time.Time) interface{} {
	return serializableAlias{
		Type:         "alias",
		MessageId:    makeMessageId(msg.MessageId, msgid),
		PreviousId:   msg.PreviousId,
		UserId:       msg.UserId,
		Timestamp:    formatTime(msg.Timestamp),
		Context:      msg.Context,
		Integrations: msg.Integrations,
	}
}

type serializableAlias struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	PreviousId   string                 `json:"previousId"`
	UserId       string                 `json:"userId"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
