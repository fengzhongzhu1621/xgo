package forwarder

import (
	"encoding/json"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/pkg/errors"
)

// messageEnvelope wraps Watermill message and contains destination topic.
type messageEnvelope struct {
	DestinationTopic string `json:"destination_topic"`

	UUID     string            `json:"uuid"`
	Payload  []byte            `json:"payload"`
	Metadata map[string]string `json:"metadata"`
}

func newMessageEnvelope(destTopic string, msg *message.Message) (*messageEnvelope, error) {
	e := &messageEnvelope{
		DestinationTopic: destTopic,
		UUID:             msg.UUID,
		Payload:          msg.Payload,
		Metadata:         msg.Metadata,
	}

	if err := e.validate(); err != nil {
		return nil, errors.Wrap(err, "cannot create a message envelope")
	}

	return e, nil
}

func (e *messageEnvelope) validate() error {
	if e.DestinationTopic == "" {
		return errors.New("unknown destination topic")
	}

	return nil
}

// wrapMessageInEnvelope 将消息放进信封
func wrapMessageInEnvelope(destinationTopic string, msg *message.Message) (*message.Message, error) {
	// 创建一个信封
	envelope, err := newMessageEnvelope(destinationTopic, msg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot envelope a message")
	}

	// 将对象转换为字节数组
	envelopedMessage, err := json.Marshal(envelope)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal a message")
	}

	// 将信封转化为一个信封消息
	wrappedMsg := message.NewMessage(uuid.NewUUID(), envelopedMessage)
	wrappedMsg.SetContext(msg.Context())

	return wrappedMsg, nil
}

func unwrapMessageFromEnvelope(msg *message.Message) (destinationTopic string, unwrappedMsg *message.Message, err error) {
	envelopedMsg := messageEnvelope{}
	if err := json.Unmarshal(msg.Payload, &envelopedMsg); err != nil {
		return "", nil, errors.Wrap(err, "cannot unmarshal message wrapped in an envelope")
	}

	if err := envelopedMsg.validate(); err != nil {
		return "", nil, errors.Wrap(err, "an unmarshalled message envelope is invalid")
	}

	watermillMessage := message.NewMessage(envelopedMsg.UUID, envelopedMsg.Payload)
	watermillMessage.Metadata = envelopedMsg.Metadata

	return envelopedMsg.DestinationTopic, watermillMessage, nil
}
