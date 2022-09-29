package publisher

import (
	. "github.com/fengzhongzhu1621/xgo/amqp/message"
)

var _ Publisher = (*DisabledPublisher)(nil)

type DisabledPublisher struct{}

func (DisabledPublisher) Publish(topic string, messages ...*Message) error {
	return ErrOutputInNoPublisherHandler
}

func (DisabledPublisher) Close() error {
	return nil
}
