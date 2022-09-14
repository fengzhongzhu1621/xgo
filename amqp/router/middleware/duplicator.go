package middleware

import (
	"xgo/amqp/message"
	"xgo/amqp/router"
)

// Duplicator is processing messages twice, to ensure that the endpoint is idempotent.
func Duplicator(h router.HandlerFunc) router.HandlerFunc {
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
