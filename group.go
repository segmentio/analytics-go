package analytics

import "time"

// This type represents object sent in a group call as described in
// https://segment.com/docs/libraries/http/#group
type Group struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	GroupId      string
	Timestamp    time.Time
	Traits       map[string]interface{}
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Group) validate() error {
	if len(msg.GroupId) == 0 {
		return FieldError{
			Type:  "analytics.Group",
			Name:  "GroupId",
			Value: msg.GroupId,
		}
	}

	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Group",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}

func (msg Group) serializable(msgid string, time time.Time) interface{} {
	return serializableGroup{
		Type:         "group",
		MessageId:    makeMessageId(msg.MessageId, msgid),
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		GroupId:      msg.GroupId,
		Timestamp:    formatTime(makeTime(msg.Timestamp, time)),
		Traits:       msg.Traits,
		Context:      msg.Context,
		Integrations: msg.Integrations,
	}
}

type serializableGroup struct {
	Type         string                 `json:"type,omitempty"`
	MessageId    string                 `json:"messageId,omitempty"`
	AnonymousId  string                 `json:"anonymousId,omitempty"`
	UserId       string                 `json:"userId,omitempty"`
	GroupId      string                 `json:"groupId"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
