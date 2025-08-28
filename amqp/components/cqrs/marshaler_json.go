package cqrs

import (
	"encoding/json"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/buildin"
	"github.com/fengzhongzhu1621/xgo/crypto/uuid"
)

var _ CommandEventMarshaler = (*JSONMarshaler)(nil)

type JSONMarshaler struct {
	NewUUID      func() string
	GenerateName func(v interface{}) string
}

func (m JSONMarshaler) Marshal(v interface{}) (*message.Message, error) {
	// 将对象转换为json字符串
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(
		m.newUUID(),
		b,
	)
	// 将对象的名称设置为消息的元信息
	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (m JSONMarshaler) newUUID() string {
	if m.NewUUID != nil {
		return m.NewUUID()
	}

	// default
	return uuid.NewUUID4()
}

func (JSONMarshaler) Unmarshal(msg *message.Message, v interface{}) (err error) {
	return json.Unmarshal(msg.Payload, v)
}

// Name 获得对象的名称
func (m JSONMarshaler) Name(cmdOrEvent interface{}) string {
	if m.GenerateName != nil {
		return m.GenerateName(cmdOrEvent)
	}

	return buildin.FullyQualifiedStructName(cmdOrEvent)
}

func (m JSONMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
