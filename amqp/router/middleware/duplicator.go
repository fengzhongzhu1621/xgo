package middleware

import (
	"github.com/fengzhongzhu1621/xgo/amqp/message"
)

// Duplicator is processing messages twice, to ensure that the endpoint is idempotent.
func Duplicator(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		firstProducedMessages, firstErr := h(msg)
		if firstErr != nil {
			return nil, firstErr
		}

		secondProducedMessages, secondErr := h(msg)
		if secondErr != nil {
			return nil, secondErr
		}

		// 消息处理两次
		return append(firstProducedMessages, secondProducedMessages...), nil
	}
}
