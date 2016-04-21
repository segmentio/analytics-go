package analytics

import "time"

// This type represents object sent in a alias call as described in
// https://segment.com/docs/libraries/http/#alias
type Alias struct {
	MessageId  string
	PreviousId string
	UserId     string
	SentAt     time.Time
	Timestamp  time.Time
}

const (
	alias = "alias"
)

func (msg Alias) validate() error {
	if len(msg.UserId) == 0 {
		return MissingFieldError{Type: alias, Name: "UserId"}
	}

	if len(msg.PreviousId) == 0 {
		return MissingFieldError{Type: alias, Name: "PreviousId"}
	}

	return nil
}

func (msg Alias) serializable() interface{} {
	return serializableAlias{
		Type:       alias,
		MessageId:  msg.MessageId,
		PreviousId: msg.PreviousId,
		UserId:     msg.UserId,
		SentAt:     formatTime(msg.SentAt),
		Timestamp:  formatTime(msg.Timestamp),
	}
}

type serializableAlias struct {
	Type       string `json:"type,omitempty"`
	MessageId  string `json:"messageId,omitempty"`
	PreviousId string `json:"previousId"`
	UserId     string `json:"userId"`
	SentAt     string `json:"sentAt,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}
