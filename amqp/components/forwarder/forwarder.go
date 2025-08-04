package forwarder

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/pkg/errors"
)

// Forwarder subscribes to the topic provided in the config and publishes them to the destination topic embedded in the enveloped message.
// 从router.Router派生的子类
type Forwarder struct {
	router    *router.Router
	publisher message.Publisher
	logger    logging.LoggerAdapter
	config    Config
}

// NewForwarder creates a forwarder which will subscribe to the topic provided in the config using the provided subscriber.
// It will publish messages received on this subscription to the destination topic embedded in the enveloped message using the provided publisher.
//
// Provided subscriber and publisher can be from different Watermill Pub/Sub implementations, i.e. MySQL subscriber and Google Pub/Sub publisher.
//
// Note: Keep in mind that by default the forwarder will nack all messages which weren't sent using a decorated publisher.
// You can change this behavior by passing a middleware which will ack them instead.
func NewForwarder(
	subscriberIn message.Subscriber,
	publisherOut message.Publisher,
	logger logging.LoggerAdapter,
	config Config,
) (*Forwarder, error) {
	config.SetDefaults()

	routerConfig := router.RouterConfig{CloseTimeout: config.CloseTimeout}
	if err := routerConfig.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid router config")
	}

	// 初始化父类
	router, err := router.NewRouter(routerConfig, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a router")
	}

	// 初始化子类
	f := &Forwarder{router, publisherOut, logger, config}

	// 添加消费者
	router.AddNoPublisherHandler(
		"events_forwarder",    // handler name
		config.ForwarderTopic, // subscribeTopic
		subscriberIn,          // message.Subscriber
		f.forwardMessage,      // NoPublishHandlerFunc
	)

	router.AddMiddleware(config.Middlewares...)

	return f, nil
}

// Run runs forwarder's handler responsible for forwarding messages.
// This call is blocking while the forwarder is running.
// ctx will be propagated to the forwarder's subscription.
//
// To stop Run() you should call Close() on the forwarder.
func (f *Forwarder) Run(ctx context.Context) error {
	return f.router.Run(ctx)
}

// Close stops forwarder's handler.
func (f *Forwarder) Close() error {
	return f.router.Close()
}

// Running returns channel which is closed when the forwarder is running.
func (f *Forwarder) Running() chan struct{} {
	return f.router.Running()
}

func (f *Forwarder) forwardMessage(msg *message.Message) error {
	// 将消息从信封拿出来
	destTopic, unwrappedMsg, err := unwrapMessageFromEnvelope(msg)
	if err != nil {
		f.logger.Error("Could not unwrap a message from an envelope", err, logging.LogFields{
			"uuid":     msg.UUID,
			"payload":  msg.Payload,
			"metadata": msg.Metadata,
			"acked":    f.config.AckWhenCannotUnwrap,
		})

		if f.config.AckWhenCannotUnwrap {
			// 忽略解码失败
			return nil
		}
		return errors.Wrap(err, "cannot unwrap message from an envelope")
	}

	if err := f.publisher.Publish(destTopic, unwrappedMsg); err != nil {
		return errors.Wrap(err, "cannot publish a message")
	}

	return nil
}
