package middleware

import (
	"xgo/amqp/message"
	"xgo/amqp/router"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// ErrInvalidPoisonQueueTopic occurs when the topic supplied to the PoisonQueue constructor is invalid.
var ErrInvalidPoisonQueueTopic = errors.New("invalid poison queue topic")

// Metadata keys which marks the reason and context why the message was deemed poisoned.
const (
	ReasonForPoisonedKey  = "reason_poisoned"
	PoisonedTopicKey      = "topic_poisoned"
	PoisonedHandlerKey    = "handler_poisoned"
	PoisonedSubscriberKey = "subscriber_poisoned"
)

// 病毒邮件队列
// 病毒邮件队列也通常为空且不可见，它里边儿保存（隔离）了一些被视为是病毒的邮件（废话），
// 这些邮件通常会导致传输服务崩溃。病毒邮件队列里的邮件不会自动尝试重新提交，管理员需要手动去删除或者恢复它们。
type poisonQueue struct {
	topic string
	pub   message.Publisher

	shouldGoToPoisonQueue func(err error) bool
}

// PoisonQueue provides a middleware that salvages unprocessable messages and published them on a separate topic.
// The main middleware chain then continues on, business as usual.
func PoisonQueue(pub message.Publisher, topic string) (router.HandlerMiddleware, error) {
	if topic == "" {
		return nil, ErrInvalidPoisonQueueTopic
	}

	pq := poisonQueue{
		topic: topic,
		pub:   pub,
		shouldGoToPoisonQueue: func(err error) bool {
			return true
		},
	}

	return pq.Middleware, nil
}

// PoisonQueueWithFilter is just like PoisonQueue, but accepts a function that decides which errors qualify for the poison queue.
func PoisonQueueWithFilter(pub message.Publisher, topic string, shouldGoToPoisonQueue func(err error) bool) (router.HandlerMiddleware, error) {
	if topic == "" {
		return nil, ErrInvalidPoisonQueueTopic
	}

	pq := poisonQueue{
		topic: topic,
		pub:   pub,

		shouldGoToPoisonQueue: shouldGoToPoisonQueue,
	}

	return pq.Middleware, nil
}

func (pq poisonQueue) publishPoisonMessage(msg *message.Message, err error) error {
	// no problems encountered, carry on
	if err == nil {
		return nil
	}

	// add context why it was poisoned
	// 将上下文的数据放到元数据中
	msg.Metadata.Set(ReasonForPoisonedKey, err.Error())
	msg.Metadata.Set(PoisonedTopicKey, router.SubscribeTopicFromCtx(msg.Context()))
	msg.Metadata.Set(PoisonedHandlerKey, router.HandlerNameFromCtx(msg.Context()))
	msg.Metadata.Set(PoisonedSubscriberKey, router.SubscriberNameFromCtx(msg.Context()))

	// don't intercept error from publish. Can't help you if the publisher is down as well.
	// 将消息发给不良队列
	return pq.pub.Publish(pq.topic, msg)
}

func (pq poisonQueue) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (events []*message.Message, err error) {
		defer func() {
			if err != nil {
				// 判断是否需要将详细发送到不良队列
				if !pq.shouldGoToPoisonQueue(err) {
					return
				}

				// handler didn't cope with the message; publish it on the poison topic and carry on as usual
				publishErr := pq.publishPoisonMessage(msg, err)
				if publishErr != nil {
					publishErr = errors.Wrap(publishErr, "cannot publish message to poison queue")
					err = multierror.Append(err, publishErr)
					return
				}

				err = nil
				return
			}
		}()

		// if h fails, the deferred function will salvage all that it can
		return h(msg)
	}
}
