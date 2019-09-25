package analytics

import "time"

var _ Message = (*Track)(nil)

// This type represents object sent in a track call as described in
// https://segment.com/docs/libraries/http/#track
type Track struct {
	// This field is exported for serialization purposes and shouldn't be set by
	// the application, its value is always overwritten by the library.
	Type string `json:"rl_type,omitempty"`

	MessageId string `json:"messageId,omitempty"`
	//AnonymousId  string       `json:"anonymousId,omitempty"`
	AnonymousId  string       `json:"rl_anonymous_id,omitempty"`
	UserId       string       `json:"rl_user_id,omitempty"`
	Event        string       `json:"rl_event"`
	Timestamp    time.Time    `json:"rl_timestamp,omitempty"`
	Context      *Context     `json:"ctx,omitempty"`
	Properties   Properties   `json:"rl_properties,omitempty"`
	Integrations Integrations `json:"integrations,omitempty"`
}

func (msg Track) internal() {
	panic(unimplementedError)
}

func (msg Track) Validate() error {
	if len(msg.Event) == 0 {
		return FieldError{
			Type:  "analytics.Track",
			Name:  "Event",
			Value: msg.Event,
		}
	}

	if len(msg.UserId) == 0 && len(msg.AnonymousId) == 0 {
		return FieldError{
			Type:  "analytics.Track",
			Name:  "UserId",
			Value: msg.UserId,
		}
	}

	return nil
}
