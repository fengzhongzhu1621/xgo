package cqrs

import (
	"github.com/fengzhongzhu1621/xgo"
	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/buildin"
	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

var _ CommandEventMarshaler = (*ProtobufMarshaler)(nil)

// ProtobufMarshaler is the default Protocol Buffers marshaler.
type ProtobufMarshaler struct {
	NewUUID      func() string
	GenerateName func(v interface{}) string
}

// Marshal marshals the given protobuf's message into watermill's Message.
func (m ProtobufMarshaler) Marshal(v interface{}) (*message.Message, error) {
	protoMsg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.WithStack(xgo.NoProtoMessageError{V: v})
	}

	b, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(
		m.newUUID(),
		b,
	)
	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (m ProtobufMarshaler) newUUID() string {
	if m.NewUUID != nil {
		return m.NewUUID()
	}

	// default
	return uuid.NewUUID()
}

// Unmarshal unmarshals given watermill's Message into protobuf's message.
func (ProtobufMarshaler) Unmarshal(msg *message.Message, v interface{}) (err error) {
	return proto.Unmarshal(msg.Payload, v.(proto.Message))
}

// Name returns the command or event's name.
func (m ProtobufMarshaler) Name(cmdOrEvent interface{}) string {
	if m.GenerateName != nil {
		return m.GenerateName(cmdOrEvent)
	}

	return buildin.FullyQualifiedStructName(cmdOrEvent)
}

// NameFromMessage returns the metadata name value for a given Message.
func (m ProtobufMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
