package analytics

import "time"

// This structure represents objects sent to the /v1/batch endpoint. We don't
// export this type because it's only meant to be used internally to send groups
// of messages in one API call.
//
// Because it implements the `Message` interface, making it public would also
// mean that programs could construct batches that embeds other batches, making
// it an invalid object construct.
// We could solve this by doing deep inspection of the `Messages` field but this
// would then have runtime costs for something that we easily solve at compile
// time.
type batch struct {
	MessageId string
	SentAt    time.Time
	Messages  []interface{}
	Context   map[string]interface{}
}

func (msg batch) validate() error {
	return nil
}

func (msg batch) serializable(msgid string, time time.Time) interface{} {
	return serializableBatch{
		MessageId: makeMessageId(msg.MessageId, msgid),
		SentAt:    formatTime(makeTime(msg.SentAt, time)),
		Messages:  msg.Messages,
		Context:   msg.Context,
	}
}

type serializableBatch struct {
	Type      string                 `json:"type,omitempty"`
	MessageId string                 `json:"messageId,omitempty"`
	SentAt    string                 `json:"sentAt,omitempty"`
	Messages  []interface{}          `json:"batch"`
	Context   map[string]interface{} `json:"context,omitempty"`
}
