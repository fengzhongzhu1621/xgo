package cqrs

import (
	"github.com/pkg/errors"
)

// Facade is a facade for creating the Command and Event buses and processors.
// It was created to avoid boilerplate, when using CQRS in the standard way.
// You can also create buses and processors manually, drawing inspiration from how it's done in NewFacade.
type Facade struct {
	commandsTopic func(commandName string) string
	commandBus    *CommandBus

	eventsTopic func(eventName string) string
	eventBus    *EventBus

	commandEventMarshaler CommandEventMarshaler
}

func (f Facade) CommandBus() *CommandBus {
	return f.commandBus
}

func (f Facade) EventBus() *EventBus {
	return f.eventBus
}

func (f Facade) CommandEventMarshaler() CommandEventMarshaler {
	return f.commandEventMarshaler
}

func NewFacade(config FacadeConfig) (*Facade, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config")
	}

	c := &Facade{
		commandsTopic:         config.GenerateCommandsTopic,
		eventsTopic:           config.GenerateEventsTopic,
		commandEventMarshaler: config.CommandEventMarshaler,
	}

	if config.CommandsEnabled() {
		var err error
		c.commandBus, err = NewCommandBus(
			config.CommandsPublisher,
			config.GenerateCommandsTopic,
			config.CommandEventMarshaler,
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create command bus")
		}
	} else {
		config.Logger.Info("Empty GenerateCommandsTopic, command bus will be not created", nil)
	}
	if config.EventsEnabled() {
		var err error
		c.eventBus, err = NewEventBus(
			config.EventsPublisher,
			config.GenerateEventsTopic,
			config.CommandEventMarshaler,
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create event bus")
		}
	} else {
		config.Logger.Info("Empty GenerateEventsTopic, event bus will be not created", nil)
	}

	if config.CommandHandlers != nil {
		commandProcessor, err := NewCommandProcessor(
			config.CommandHandlers(c.commandBus, c.eventBus),
			config.GenerateCommandsTopic,
			config.CommandsSubscriberConstructor,
			config.CommandEventMarshaler,
			config.Logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create command processor")
		}

		if err := commandProcessor.AddHandlersToRouter(config.Router); err != nil {
			return nil, err
		}
	}
	if config.EventHandlers != nil {
		eventProcessor, err := NewEventProcessor(
			config.EventHandlers(c.commandBus, c.eventBus),
			config.GenerateEventsTopic,
			config.EventsSubscriberConstructor,
			config.CommandEventMarshaler,
			config.Logger,
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create event processor")
		}

		if err := eventProcessor.AddHandlersToRouter(config.Router); err != nil {
			return nil, err
		}
	}

	return c, nil
}
