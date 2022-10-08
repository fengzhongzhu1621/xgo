package fanin

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/pkg/errors"
)

var _ router.IRouterConfig = (*Config)(nil)

type Config struct {
	// SourceTopics contains topics on which FanIn subscribes.
	SourceTopics []string // 用于一个router监听多个topic

	// TargetTopic determines the topic on which messages from SourceTopics are published.
	TargetTopic string

	// CloseTimeout determines how long router should work for handlers when closing.
	CloseTimeout time.Duration
}

func (c *Config) SetDefaults() {
	if c.CloseTimeout == 0 {
		c.CloseTimeout = time.Second * 30
	}
}

func (c *Config) Validate() error {
	if len(c.SourceTopics) == 0 {
		return errors.New("sourceTopics must not be empty")
	}

	for _, fromTopic := range c.SourceTopics {
		if fromTopic == "" {
			return errors.New("sourceTopics must not be empty")
		}
	}

	if c.TargetTopic == "" {
		return errors.New("targetTopic must not be empty")
	}

	for _, fromTopic := range c.SourceTopics {
		if fromTopic == c.TargetTopic {
			return errors.New("sourceTopics must not contain targetTopic")
		}
	}

	return nil
}
