package analytics

import "time"

// This type represents object sent in a group call as described in
// https://segment.com/docs/libraries/http/#group
type Group struct {
	MessageId    string
	AnonymousId  string
	UserId       string
	GroupId      string
	SentAt       time.Time
	Timestamp    time.Time
	Traits       map[string]interface{}
	Context      map[string]interface{}
	Integrations map[string]interface{}
}

func (msg Group) serializable() interface{} {
	return serializableGroup{
		Type:         "group",
		MessageId:    msg.MessageId,
		AnonymousId:  msg.AnonymousId,
		UserId:       msg.UserId,
		GroupId:      msg.GroupId,
		SentAt:       formatTime(msg.SentAt),
		Timestamp:    formatTime(msg.Timestamp),
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
	SentAt       string                 `json:"sentAt,omitempty"`
	Timestamp    string                 `json:"timestamp,omitempty"`
	Traits       map[string]interface{} `json:"traits,omitempty"`
	Context      map[string]interface{} `json:"context,omitempty"`
	Integrations map[string]interface{} `json:"integrations,omitempty"`
}
