package cqrs

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/pkg/errors"
)

// EventHandler receives events defined by NewEvent and handles them with its Handle method.
// If using DDD, CommandHandler may modify and persist the aggregate.
// It can also invoke a process manager, a saga or just build a read model.
//
// In contrast to CommandHandler, every Event can have multiple EventHandlers.
//
// One instance of EventHandler is used during handling messages.
// When multiple events are delivered at the same time, Handle method can be executed multiple times at the same time.
// Because of that, Handle method needs to be thread safe!
type EventHandler interface {
	// HandlerName is the name used in message.Router while creating handler.
	//
	// It will be also passed to EventsSubscriberConstructor.
	// May be useful, for example, to create a consumer group per each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result with **reconsuming all messages** !!!
	HandlerName() string // 路由处理器的名称

	NewEvent() interface{} // 获得消息的类型

	Handle(ctx context.Context, event interface{}) error // 消息处理函数
}

// EventsSubscriberConstructor creates a subscriber for EventHandler.
// It allows you to create separated customized Subscriber for every command handler.
// 获得handlerName获得对应的消费者对象
type EventsSubscriberConstructor func(handlerName string) (message.Subscriber, error)

// EventProcessor determines which EventHandler should handle event received from event bus.
type EventProcessor struct {
	handlers      []EventHandler // 多个处理器
	generateTopic func(eventName string) string

	subscriberConstructor EventsSubscriberConstructor // 获得handlerName获得对应的消费者对象

	marshaler CommandEventMarshaler // 消息编解码器
	logger    logging.LoggerAdapter
}

func NewEventProcessor(
	handlers []EventHandler,
	generateTopic func(eventName string) string,
	subscriberConstructor EventsSubscriberConstructor,
	marshaler CommandEventMarshaler,
	logger logging.LoggerAdapter,
) (*EventProcessor, error) {
	if len(handlers) == 0 {
		return nil, errors.New("missing handlers")
	}
	if generateTopic == nil {
		return nil, errors.New("nil generateTopic")
	}
	if subscriberConstructor == nil {
		return nil, errors.New("missing subscriberConstructor")
	}
	if marshaler == nil {
		return nil, errors.New("missing marshaler")
	}
	if logger == nil {
		logger = logging.NopLogger{}
	}

	return &EventProcessor{
		handlers,
		generateTopic,
		subscriberConstructor,
		marshaler,
		logger,
	}, nil
}

func (p EventProcessor) AddHandlersToRouter(r *router.Router) error {
	for i := range p.Handlers() {
		handler := p.handlers[i]
		// 获得路由handler的名称
		handlerName := handler.HandlerName()
		// 获得消息的类型
		eventName := p.marshaler.Name(handler.NewEvent())
		// 根据消息类型生成topic
		topicName := p.generateTopic(eventName)

		logger := p.logger.With(logging.LogFields{
			"event_handler_name": handlerName,
			"topic":              topicName,
		})

		// 消息处理函数
		handlerFunc, err := p.routerHandlerFunc(handler, logger)
		if err != nil {
			return err
		}

		logger.Debug("Adding CQRS event handler to router", nil)

		// 获得handlerName获得对应的消费者对象
		subscriber, err := p.subscriberConstructor(handlerName)
		if err != nil {
			return errors.Wrap(err, "cannot create subscriber for event processor")
		}
		// 将消息处理器附加到router上
		r.AddNoPublisherHandler(
			handlerName,
			topicName,
			subscriber,
			handlerFunc,
		)
	}

	return nil
}

func (p EventProcessor) Handlers() []EventHandler {
	return p.handlers
}

// routerHandlerFunc 构造消费者的消息处理函数装饰器
func (p EventProcessor) routerHandlerFunc(
	handler EventHandler,
	logger logging.LoggerAdapter,
) (router.NoPublishHandlerFunc, error) {
	initEvent := handler.NewEvent()
	if err := p.validateEvent(initEvent); err != nil {
		return nil, err
	}

	expectedEventName := p.marshaler.Name(initEvent)

	return func(msg *message.Message) error {
		// 获得消息的类型
		event := handler.NewEvent()
		messageEventName := p.marshaler.NameFromMessage(msg)

		// 判断消息类型是否一致
		if messageEventName != expectedEventName {
			logger.Trace("Received different event type than expected, ignoring", logging.LogFields{
				"message_uuid":        msg.UUID,
				"expected_event_type": expectedEventName,
				"received_event_type": messageEventName,
			})
			return nil
		}

		logger.Debug("Handling event", logging.LogFields{
			"message_uuid":        msg.UUID,
			"received_event_type": messageEventName,
		})

		// 解码消息
		if err := p.marshaler.Unmarshal(msg, event); err != nil {
			return err
		}

		// 处理消息
		if err := handler.Handle(msg.Context(), event); err != nil {
			logger.Debug("Error when handling event", logging.LogFields{"err": err})
			return err
		}

		return nil
	}, nil
}

func (p EventProcessor) validateEvent(event interface{}) error {
	// EventHandler's NewEvent must return a pointer, because it is used to unmarshal
	if err := reflectutils.IsPointer(event); err != nil {
		return errors.Wrap(err, "command must be a non-nil pointer")
	}

	return nil
}
