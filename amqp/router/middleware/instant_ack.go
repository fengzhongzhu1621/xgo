package middleware

import (
	"xgo/amqp/message"
	"xgo/amqp/router"
)

// InstantAck makes the handler instantly acknowledge the incoming message, regardless of any errors.
// It may be used to gain throughput, but at a cost:
// If you had exactly-once delivery, you may expect at-least-once instead.
// If you had ordered messages, the ordering might be broken.
func InstantAck(h router.HandlerFunc) router.HandlerFunc {
	return func(message *message.Message) ([]*message.Message, error) {
		// 在处理消息前进行ack
		message.Ack()
		return h(message)
	}
}
