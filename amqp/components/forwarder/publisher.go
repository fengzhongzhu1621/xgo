package forwarder

import (
	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/pkg/errors"
)

var _ router.IRouterConfig = (*PublisherConfig)(nil)

type PublisherConfig struct {
	// ForwarderTopic is a topic which the forwarder is listening to. Publisher will send enveloped messages to this topic.
	// Defaults to `forwarder_topic`.
	ForwarderTopic string
}

func (c *PublisherConfig) SetDefaults() {
	if c.ForwarderTopic == "" {
		c.ForwarderTopic = defaultForwarderTopic
	}
}

func (c *PublisherConfig) Validate() error {
	if c.ForwarderTopic == "" {
		return errors.New("empty forwarder topic")
	}

	return nil
}

var _ message.Publisher = (*Publisher)(nil)

// Publisher changes `Publish` method behavior so it wraps a sent message in an envelope
// and sends it to the forwarder topic provided in the config.
type Publisher struct {
	wrappedPublisher message.Publisher
	config           PublisherConfig
}

func NewPublisher(publisher message.Publisher, config PublisherConfig) *Publisher {
	config.SetDefaults()

	return &Publisher{
		wrappedPublisher: publisher,
		config:           config,
	}
}

func (p *Publisher) Publish(topic string, messages ...*message.Message) error {
	envelopedMessages := make([]*message.Message, 0, len(messages))
	for _, msg := range messages {
		// 将消息装进信封，消息中的topic和发送topic不一样的
		envelopedMsg, err := wrapMessageInEnvelope(topic, msg)
		if err != nil {
			return errors.Wrapf(
				err,
				"cannot wrap message, target topic: '%s', uuid: '%s'",
				topic,
				msg.UUID,
			)
		}

		envelopedMessages = append(envelopedMessages, envelopedMsg)
	}

	if err := p.wrappedPublisher.Publish(p.config.ForwarderTopic, envelopedMessages...); err != nil {
		return errors.Wrapf(
			err,
			"cannot publish messages to forwarder topic: '%s'",
			p.config.ForwarderTopic,
		)
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.wrappedPublisher.Close()
}
