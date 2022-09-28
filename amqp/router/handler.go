package router

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	. "xgo/amqp/message"
	. "xgo/amqp/publisher"
	"xgo/log"

	"github.com/pkg/errors"
)

type handler struct {
	name   string
	logger log.LoggerAdapter

	subscriber     Subscriber
	subscribeTopic string
	subscriberName string

	publisher     Publisher
	publishTopic  string
	publisherName string

	handlerFunc HandlerFunc

	runningHandlersWg *sync.WaitGroup // 标记所有handler接收到的一批消息是否都处理完毕

	messagesCh <-chan *Message // 从队列读取消息，将格式化后的的消息放入此队列

	started   bool
	startedCh chan struct{}

	stopFn         context.CancelFunc
	stopped        chan struct{} // 标记handler是否结束
	routersCloseCh chan struct{}
}

func (h *handler) run(ctx context.Context, middlewares []Middleware) {
	h.logger.Info("Starting handler", log.LogFields{
		"subscriber_name": h.name,
		"topic":           h.subscribeTopic,
	})

	// 消息处理函数
	middlewareHandler := h.handlerFunc

	// 对消息处理函数添加装饰器
	// first added middlewares should be executed first (so should be at the top of call stack)
	for i := len(middlewares) - 1; i >= 0; i-- {
		currentMiddleware := middlewares[i]
		// 选择匹配的消息处理函数装饰器
		isValidHandlerLevelMiddleware := currentMiddleware.HandlerName == h.name
		if currentMiddleware.IsRouterLevel || isValidHandlerLevelMiddleware {
			// 装饰消息处理函数
			middlewareHandler = currentMiddleware.Handler(middlewareHandler)
		}
	}
	// 注册处理路由关闭信号和上下文关闭信号协程
	go h.handleClose(ctx)

	// 阻塞从队列中获取一条格式化后的消息
	// todo 可优化，会创建大量的 go h.handleMessage，导致资源耗尽
	for msg := range h.messagesCh {
		// 标记正在处理一个消息
		h.runningHandlersWg.Add(1)
		// 处理消息
		go h.handleMessage(msg, middlewareHandler)
	}

	// 如果消费者销毁，则释放资源
	if h.publisher != nil {
		h.logger.Debug("Waiting for publisher to close", nil)
		if err := h.publisher.Close(); err != nil {
			h.logger.Error("Failed to close publisher", err, nil)
		}
		h.logger.Debug("Publisher closed", nil)
	}

	h.logger.Debug("Router handler stopped", nil)
	close(h.stopped)
}

// addHandlerContext enriches the contex with values that are relevant within this handler's context.
func (h *handler) addHandlerContext(messages ...*Message) {
	for i, msg := range messages {
		ctx := msg.Context()

		if h.name != "" {
			ctx = context.WithValue(ctx, handlerNameKey, h.name)
		}
		if h.publisherName != "" {
			ctx = context.WithValue(ctx, publisherNameKey, h.publisherName)
		}
		if h.subscriberName != "" {
			ctx = context.WithValue(ctx, subscriberNameKey, h.subscriberName)
		}
		if h.subscribeTopic != "" {
			ctx = context.WithValue(ctx, subscribeTopicKey, h.subscribeTopic)
		}
		if h.publishTopic != "" {
			ctx = context.WithValue(ctx, publishTopicKey, h.publishTopic)
		}
		messages[i].SetContext(ctx)
	}
}

func (h *handler) handleClose(ctx context.Context) {
	select {
	case <-h.routersCloseCh:
		// for backward compatibility we are closing subscriber
		// router关闭时，handler会接收到此信号
		h.logger.Debug("Waiting for subscriber to close", nil)
		if err := h.subscriber.Close(); err != nil {
			h.logger.Error("Failed to close subscriber", err, nil)
		}
		h.logger.Debug("Subscriber closed", nil)
	case <-ctx.Done():
		// we are closing subscriber just when entire router is closed
	}
	// 向父上下文发送关闭信号
	h.stopFn()
}

func (h *handler) handleMessage(msg *Message, handler HandlerFunc) {
	// 标记一个消息处理完成
	defer h.runningHandlersWg.Done()
	msgFields := log.LogFields{"message_uuid": msg.UUID}

	defer func() {
		if recovered := recover(); recovered != nil {
			h.logger.Error(
				"Panic recovered in handler. Stack: "+string(debug.Stack()),
				errors.Errorf("%s", recovered),
				msgFields,
			)
			msg.Nack()
		}
	}()

	h.logger.Trace("Received message", msgFields)

	// 执行消息处理函数，可能返回多个消息给下一个生产者
	// handler包含多个handler中间件，可以处理复杂的逻辑
	producedMessages, err := handler(msg)
	if err != nil {
		h.logger.Error("Handler returned error", err, nil)
		msg.Nack()
		return
	}

	// 将需要发给下一个生产者的消息添加上下文
	h.addHandlerContext(producedMessages...)

	// 发送处理后的消息给下一个生产者
	if err := h.publishProducedMessages(producedMessages, msgFields); err != nil {
		h.logger.Error("Publishing produced messages failed", err, nil)
		msg.Nack()
		return
	}

	// 消息应答
	msg.Ack()
	h.logger.Trace("Message acked", msgFields)
}

func (h *handler) publishProducedMessages(producedMessages Messages, msgFields log.LogFields) error {
	if len(producedMessages) == 0 {
		return nil
	}

	if h.publisher == nil {
		return ErrOutputInNoPublisherHandler
	}

	h.logger.Trace("Sending produced messages", msgFields.Add(log.LogFields{
		"produced_messages_count": len(producedMessages),
		"publish_topic":           h.publishTopic,
	}))

	for _, msg := range producedMessages {
		if err := h.publisher.Publish(h.publishTopic, msg); err != nil {
			// todo - how to deal with it better/transactional/retry?
			h.logger.Error("Cannot publish message", err, msgFields.Add(log.LogFields{
				"not_sent_message": fmt.Sprintf("%#v", producedMessages),
			}))

			return err
		}
	}

	return nil
}

// Handler handles Messages.
type Handler struct {
	router  *Router
	handler *handler
}

// AddMiddleware adds new middleware to the specified handler in the router.
//
// The order of middleware matters. Middleware added at the beginning is executed first.
func (h *Handler) AddMiddleware(m ...HandlerMiddleware) {
	handler := h.handler
	handler.logger.Debug("Adding middleware to handler", log.LogFields{
		"count":       fmt.Sprintf("%d", len(m)),
		"handlerName": handler.name,
	})

	h.router.addHandlerLevelMiddleware(handler.name, m...)
}

// Started returns channel which is stopped when handler is running.
func (h *Handler) Started() chan struct{} {
	return h.handler.startedCh
}

// Stop stops the handler.
// Stop is asynchronous.
// You can check if handler was stopped with Stopped() function.
func (h *Handler) Stop() {
	if !h.handler.started {
		panic("handler is not started")
	}

	h.handler.stopFn()
}

// Stopped returns channel which is stopped when handler did stop.
func (h *Handler) Stopped() chan struct{} {
	return h.handler.stopped
}
