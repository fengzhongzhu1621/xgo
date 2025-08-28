package forwarder

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/pkg/errors"
)

const defaultForwarderTopic = "forwarder_topic"

var _ router.IRouterConfig = (*Config)(nil)

type Config struct {
	// ForwarderTopic is a topic on which the forwarder will be listening to enveloped messages to forward.
	// Defaults to `forwarder_topic`.
	ForwarderTopic string

	// Middlewares are used to decorate forwarder's handler function.
	Middlewares []router.HandlerMiddleware

	// CloseTimeout determines how long router should work for handlers when closing.
	CloseTimeout time.Duration

	// AckWhenCannotUnwrap enables acking of messages which cannot be unwrapped from an envelope.
	AckWhenCannotUnwrap bool
}

func (c *Config) SetDefaults() {
	if c.CloseTimeout == 0 {
		c.CloseTimeout = time.Second * 30
	}
	if c.ForwarderTopic == "" {
		c.ForwarderTopic = defaultForwarderTopic
	}
}

func (c *Config) Validate() error {
	if c.ForwarderTopic == "" {
		return errors.New("empty forwarder topic")
	}

	return nil
}
