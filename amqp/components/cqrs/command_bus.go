package cqrs

import (
	"context"

	"github.com/pkg/errors"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
)

// CommandBus transports commands to command handlers.
// 用于将接收到的command对象分发到指定的handler
type CommandBus struct {
	publisher     message.Publisher
	generateTopic func(commandName string) string
	marshaler     CommandEventMarshaler // 消息编解码器
}

func NewCommandBus(
	publisher message.Publisher,
	generateTopic func(commandName string) string,
	marshaler CommandEventMarshaler,
) (*CommandBus, error) {
	if publisher == nil {
		return nil, errors.New("missing publisher")
	}
	if generateTopic == nil {
		return nil, errors.New("missing generateTopic")
	}
	if marshaler == nil {
		return nil, errors.New("missing marshaler")
	}

	return &CommandBus{publisher, generateTopic, marshaler}, nil
}

// Send sends command to the command bus.
func (c CommandBus) Send(ctx context.Context, cmd interface{}) error {
	// 将cmd对象转换为消息
	msg, err := c.marshaler.Marshal(cmd)
	if err != nil {
		return err
	}
	// 获得topic名称
	commandName := c.marshaler.Name(cmd)
	topicName := c.generateTopic(commandName)

	msg.SetContext(ctx)

	// 发送消息给队列
	return c.publisher.Publish(topicName, msg)
}
