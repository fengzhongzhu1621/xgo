package middleware

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
)

// Timeout makes the handler cancel the incoming message's context after a specified time.
// Any timeout-sensitive functionality of the handler should listen on msg.Context().Done() to know when to fail.
func Timeout(timeout time.Duration) func(message.HandlerFunc) message.HandlerFunc {
	return func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			// 创建一个子节点的context, timeout秒后自动超时
			ctx, cancel := context.WithTimeout(msg.Context(), timeout)
			defer func() {
				cancel()
			}()

			msg.SetContext(ctx)
			// 超时取消消息处理
			return h(msg)
		}
	}
}
