package journify

import "time"

var _ Message = (*Identify)(nil)

type Identify struct {
	// This field is exported for serialization purposes and shouldn't be set by
	// the application, its value is always overwritten by the library.
	Type string `json:"type,omitempty"`

	MessageId   string    `json:"messageId,omitempty"`
	AnonymousId string    `json:"anonymousId,omitempty"`
	UserId      string    `json:"userId,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Context     *Context  `json:"context,omitempty"`
	Traits      Traits    `json:"traits,omitempty"`
	WriteKey    string    `json:"writeKey,omitempty"`
}

func (msg Identify) Validate() error {
	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "journify.Identify",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}
