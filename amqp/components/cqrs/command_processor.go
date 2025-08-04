package cqrs

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/pkg/errors"
)

// CommandHandler receives a command defined by NewCommand and handles it with the Handle method.
// If using DDD, CommandHandler may modify and persist the aggregate.
//
// In contrast to EvenHandler, every Command must have only one CommandHandler.
//
// One instance of CommandHandler is used during handling messages.
// When multiple commands are delivered at the same time, Handle method can be executed multiple times at the same time.
// Because of that, Handle method needs to be thread safe!
type CommandHandler interface {
	// HandlerName is the name used in message.Router while creating handler.
	//
	// It will be also passed to CommandsSubscriberConstructor.
	// May be useful, for example, to create a consumer group per each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result with **reconsuming all messages**!
	HandlerName() string

	NewCommand() interface{}

	Handle(ctx context.Context, cmd interface{}) error
}

// CommandsSubscriberConstructor creates subscriber for CommandHandler.
// It allows you to create a separate customized Subscriber for every command handler.
type CommandsSubscriberConstructor func(handlerName string) (message.Subscriber, error)

// CommandProcessor determines which CommandHandler should handle the command received from the command bus.
type CommandProcessor struct {
	handlers      []CommandHandler // 命令处理器
	generateTopic func(commandName string) string

	subscriberConstructor CommandsSubscriberConstructor // 获取消费者

	marshaler CommandEventMarshaler
	logger    logging.LoggerAdapter
}

// NewCommandProcessor creates a new CommandProcessor.
func NewCommandProcessor(
	handlers []CommandHandler,
	generateTopic func(commandName string) string,
	subscriberConstructor CommandsSubscriberConstructor,
	marshaler CommandEventMarshaler,
	logger logging.LoggerAdapter,
) (*CommandProcessor, error) {
	if len(handlers) == 0 {
		return nil, errors.New("missing handlers")
	}
	if generateTopic == nil {
		return nil, errors.New("missing generateTopic")
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

	return &CommandProcessor{
		handlers,
		generateTopic,
		subscriberConstructor,
		marshaler,
		logger,
	}, nil
}

// DuplicateCommandHandlerError occurs when a handler with the same name already exists.
type DuplicateCommandHandlerError struct {
	CommandName string
}

func (d DuplicateCommandHandlerError) Error() string {
	return fmt.Sprintf("command handler for command %s already exists", d.CommandName)
}

// AddHandlersToRouter adds the CommandProcessor's handlers to the given router.
func (p CommandProcessor) AddHandlersToRouter(r *router.Router) error {
	handledCommands := map[string]struct{}{}

	// 遍历所有的命令处理器，通常一种命令只能有一个命令处理器
	for i := range p.Handlers() {
		handler := p.handlers[i]
		handlerName := handler.HandlerName()
		commandName := p.marshaler.Name(handler.NewCommand())
		topicName := p.generateTopic(commandName)
		// 忽略重复的命令
		if _, ok := handledCommands[commandName]; ok {
			return DuplicateCommandHandlerError{commandName}
		}
		handledCommands[commandName] = struct{}{}

		logger := p.logger.With(logging.LogFields{
			"command_handler_name": handlerName,
			"topic":                topicName,
		})
		// 生成消息处理函数
		handlerFunc, err := p.routerHandlerFunc(handler, logger)
		if err != nil {
			return err
		}

		logger.Debug("Adding CQRS command handler to router", nil)
		// 获得消费者，通常每个命令处理器对应一个消费者
		subscriber, err := p.subscriberConstructor(handlerName)
		if err != nil {
			return errors.Wrap(err, "cannot create subscriber for command processor")
		}
		// 将消息处理器添加到router上
		r.AddNoPublisherHandler(
			handlerName,
			topicName,
			subscriber,
			handlerFunc,
		)
	}

	return nil
}

// Handlers returns the CommandProcessor's handlers.
func (p CommandProcessor) Handlers() []CommandHandler {
	return p.handlers
}

// routerHandlerFunc 构造消息处理器装饰器
func (p CommandProcessor) routerHandlerFunc(
	handler CommandHandler,
	logger logging.LoggerAdapter,
) (router.NoPublishHandlerFunc, error) {
	cmd := handler.NewCommand()
	cmdName := p.marshaler.Name(cmd)

	if err := p.validateCommand(cmd); err != nil {
		return nil, err
	}

	return func(msg *message.Message) error {
		cmd := handler.NewCommand()
		messageCmdName := p.marshaler.NameFromMessage(msg)

		if messageCmdName != cmdName {
			logger.Trace(
				"Received different command type than expected, ignoring",
				logging.LogFields{
					"message_uuid":          msg.UUID,
					"expected_command_type": cmdName,
					"received_command_type": messageCmdName,
				},
			)
			return nil
		}

		logger.Debug("Handling command", logging.LogFields{
			"message_uuid":          msg.UUID,
			"received_command_type": messageCmdName,
		})

		if err := p.marshaler.Unmarshal(msg, cmd); err != nil {
			return err
		}

		if err := handler.Handle(msg.Context(), cmd); err != nil {
			logger.Debug("Error when handling command", logging.LogFields{"err": err})
			return err
		}

		return nil
	}, nil
}

func (p CommandProcessor) validateCommand(cmd interface{}) error {
	// CommandHandler's NewCommand must return a pointer, because it is used to unmarshal
	if err := reflectutils.IsPointer(cmd); err != nil {
		return errors.Wrap(err, "command must be a non-nil pointer")
	}

	return nil
}
