package router

import (
	. "github.com/fengzhongzhu1621/xgo/amqp/message"
)

// HandlerMiddleware allows us to write something like decorators to HandlerFunc.
// It can execute something before handler (for example: modify consumed message)
// or after (modify produced messages, ack/nack on consumed message, handle errors, logging, etc.).
//
// It can be attached to the router by using `AddMiddleware` method.
//
// Example:
//		func ExampleMiddleware(h message.HandlerFunc) message.HandlerFunc {
//			return func(message *message.Message) ([]*message.Message, error) {
//				fmt.Println("executed before handler")
//				producedMessages, err := h(message)
//				fmt.Println("executed after handler")
//
//				return producedMessages, err
//			}
//		}
// 消息中间件装饰器
type HandlerMiddleware func(h HandlerFunc) HandlerFunc

type Middleware struct {
	Handler       HandlerMiddleware // 消息处理函数，使用中间件装饰器进行装饰
	HandlerName   string
	IsRouterLevel bool // 是否用在router中
}
