package gochannel

import (
	"context"
	"sync"

	"github.com/fengzhongzhu1621/xgo/crypto/uuid"

	"github.com/fengzhongzhu1621/xgo/logging"

	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
)

// Config holds the GoChannel Pub/Sub's configuration options.
type Config struct {
	// Output channel buffer size.
	OutputChannelBuffer int64 // 订阅者缓存队列的大小

	// If persistent is set to true, when subscriber subscribes to the topic,
	// it will receive all previously produced messages.
	//
	// All messages are persisted to the memory (simple slice),
	// so be aware that with large amount of messages you can go out of the memory.
	Persistent bool // 是否持久化存储

	// When true, Publish will block until subscriber Ack's the message.
	// If there are no subscribers, Publish will not block (also when Persistent is true).
	BlockPublishUntilSubscriberAck bool // 是否等待所有的订阅者都处理完同一个消息
}

// GoChannel is the simplest Pub/Sub implementation.
// It is based on Golang's channels which are sent within the process.
//
// GoChannel has no global state,
// that means that you need to use the same instance for Publishing and Subscribing!
//
// When GoChannel is persistent, messages order is not guaranteed.
type GoChannel struct {
	config Config // 生产者配置
	logger logging.LoggerAdapter

	subscribersWg          sync.WaitGroup
	subscribers            map[string][]*GoChannelSubscriber // 多个消费者
	subscribersLock        sync.RWMutex
	subscribersByTopicLock sync.Map // map of *sync.Mutex

	closed     bool // 关闭标识
	closedLock sync.Mutex
	closing    chan struct{}

	persistedMessages     map[string][]*message.Message
	persistedMessagesLock sync.RWMutex
}

// NewGoChannel creates new GoChannel Pub/Sub.
//
// This GoChannel is not persistent.
// That means if you send a message to a topic to which no subscriber is subscribed, that message will be discarded.
func NewGoChannel(config Config, logger logging.LoggerAdapter) *GoChannel {
	if logger == nil {
		logger = logging.NopLogger{}
	}

	return &GoChannel{
		config: config,

		subscribers:            make(map[string][]*GoChannelSubscriber), // 多个订阅者
		subscribersByTopicLock: sync.Map{},
		logger: logger.With(logging.LogFields{
			"pubsub_uuid": shortuuid.New(),
		}),

		closing: make(chan struct{}),

		persistedMessages: map[string][]*message.Message{}, // 存放持久化的消息
	}
}

// Publish in GoChannel is NOT blocking until all consumers consume.
// Messages will be send in background.
//
// Messages may be persisted or not, depending of persistent attribute.
func (g *GoChannel) Publish(topic string, messages ...*message.Message) error {
	if g.isClosed() {
		return errors.New("Pub/Sub closed")
	}

	// 复制消息
	messagesToPublish := make(message.Messages, len(messages))
	for i, msg := range messages {
		messagesToPublish[i] = msg.Copy()
	}

	g.subscribersLock.RLock()
	defer g.subscribersLock.RUnlock()

	subLock, _ := g.subscribersByTopicLock.LoadOrStore(topic, &sync.Mutex{})
	subLock.(*sync.Mutex).Lock()
	defer subLock.(*sync.Mutex).Unlock()

	// 发送消息后，将消息持久化存储，防止丢失
	// 1. 在没有消费者的情况下，先缓冲生产者发送的消息
	// 2. 在消费者启动后，再接收缓冲的消息
	if g.config.Persistent {
		g.persistedMessagesLock.Lock()
		if _, ok := g.persistedMessages[topic]; !ok {
			g.persistedMessages[topic] = make([]*message.Message, 0)
		}
		g.persistedMessages[topic] = append(g.persistedMessages[topic], messagesToPublish...)
		g.persistedMessagesLock.Unlock()
	}

	// 发送消息给所有的订阅者
	for i := range messagesToPublish {
		msg := messagesToPublish[i]

		// 发送消息给订阅者，返回所有的订阅者都处理了消息的信号
		ackedBySubscribers, err := g.sendMessage(topic, msg)
		if err != nil {
			return err
		}

		// 等待所有的订阅者处理完成，默认不等待
		if g.config.BlockPublishUntilSubscriberAck {
			g.waitForAckFromSubscribers(msg, ackedBySubscribers)
		}
	}

	return nil
}

// waitForAckFromSubscribers 等待订阅者返回ack消息
func (g *GoChannel) waitForAckFromSubscribers(msg *message.Message, ackedByConsumer <-chan struct{}) {
	logFields := logging.LogFields{"message_uuid": msg.UUID}
	g.logger.Debug("Waiting for subscribers ack", logFields)

	select {
	case <-ackedByConsumer:
		// 等待所有的订阅者处理完毕
		g.logger.Trace("Message acked by subscribers", logFields)
	case <-g.closing:
		//Gochannl关闭
		g.logger.Trace("Closing Pub/Sub before ack from subscribers", logFields)
	}
}

// 生产者发送消息给所有注册的订阅者，并启动协程模拟订阅者接受消息到消息缓存队列，并等待消息返回ack
func (g *GoChannel) sendMessage(topic string, message *message.Message) (<-chan struct{}, error) {
	// 根据topic选择多个订阅者
	subscribers := g.topicSubscribers(topic)
	ackedBySubscribers := make(chan struct{})

	logFields := logging.LogFields{"message_uuid": message.UUID, "topic": topic}

	if len(subscribers) == 0 {
		close(ackedBySubscribers)
		g.logger.Info("No subscribers to send message", logFields)
		return ackedBySubscribers, nil
	}

	// 多个消费者处理同一个消息
	go func(subscribers []*GoChannelSubscriber) {
		wg := &sync.WaitGroup{}

		for i := range subscribers {
			subscriber := subscribers[i]

			wg.Add(1)
			go func() {
				// 模拟订阅者接受消息
				subscriber.SendMessageToSubscriber(message, logFields)
				wg.Done()
			}()
		}
		// 等待所有的订阅者都处理了这个消息
		wg.Wait()
		close(ackedBySubscribers)
	}(subscribers)

	return ackedBySubscribers, nil
}

// Subscribe returns channel to which all published messages are sent.
// Messages are not persisted. If there are no subscribers and message is produced it will be gone.
//
// There are no consumer groups support etc. Every consumer will receive every produced message.
// 创建一个订阅者并返回接收消息的缓存队列
func (g *GoChannel) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	g.closedLock.Lock()

	if g.closed {
		return nil, errors.New("Pub/Sub closed")
	}

	g.subscribersWg.Add(1)
	g.closedLock.Unlock()

	g.subscribersLock.Lock()

	subLock, _ := g.subscribersByTopicLock.LoadOrStore(topic, &sync.Mutex{})
	subLock.(*sync.Mutex).Lock()

	// 自动创建一个订阅者
	s := &GoChannelSubscriber{
		ctx:           ctx,
		uuid:          uuid.NewUUID(),
		outputChannel: make(chan *message.Message, g.config.OutputChannelBuffer),
		logger:        g.logger,
		closing:       make(chan struct{}),
	}

	// 订阅者注销协程
	go func(s *GoChannelSubscriber, g *GoChannel) {
		select {
		case <-ctx.Done():
			// unblock
		case <-g.closing:
			// unblock
		}

		// 生产者关闭时，需要关闭订阅者
		s.Close()

		g.subscribersLock.Lock()
		defer g.subscribersLock.Unlock()

		subLock, _ := g.subscribersByTopicLock.Load(topic)
		subLock.(*sync.Mutex).Lock()
		defer subLock.(*sync.Mutex).Unlock()

		// 注销订阅者
		g.removeSubscriber(topic, s)
		g.subscribersWg.Done()
	}(s, g)

	if !g.config.Persistent {
		defer g.subscribersLock.Unlock()
		defer subLock.(*sync.Mutex).Unlock()

		// 注册消费者
		g.addSubscriber(topic, s)

		return s.outputChannel, nil
	}

	go func(s *GoChannelSubscriber) {
		defer g.subscribersLock.Unlock()
		defer subLock.(*sync.Mutex).Unlock()

		// 获取持久化的消息
		g.persistedMessagesLock.RLock()
		messages, ok := g.persistedMessages[topic]
		g.persistedMessagesLock.RUnlock()
		// 将持久化的消息发给订阅者
		// 因为开始发送消息时，订阅者可能不存在，这时候的消息需要缓存下来
		// 等到订阅者注册后，再从缓存中把重新发给订阅者
		if ok {
			for i := range messages {
				msg := g.persistedMessages[topic][i]
				logFields := logging.LogFields{"message_uuid": msg.UUID, "topic": topic}

				go s.SendMessageToSubscriber(msg, logFields)
			}
		}

		g.addSubscriber(topic, s)
	}(s)

	return s.outputChannel, nil
}

func (g *GoChannel) addSubscriber(topic string, s *GoChannelSubscriber) {
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*GoChannelSubscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], s)
}

func (g *GoChannel) removeSubscriber(topic string, toRemove *GoChannelSubscriber) {
	removed := false
	for i, sub := range g.subscribers[topic] {
		if sub == toRemove {
			g.subscribers[topic] = append(g.subscribers[topic][:i], g.subscribers[topic][i+1:]...)
			removed = true
			break
		}
	}
	if !removed {
		panic("cannot remove subscriber, not found " + toRemove.uuid)
	}
}

// topicSubscribers 选择多个订阅者
func (g *GoChannel) topicSubscribers(topic string) []*GoChannelSubscriber {
	subscribers, ok := g.subscribers[topic]
	if !ok {
		return nil
	}

	// let's do a copy to avoid race conditions and deadlocks due to lock
	// 复制订阅者
	subscribersCopy := make([]*GoChannelSubscriber, len(subscribers))
	for i, s := range subscribers {
		subscribersCopy[i] = s
	}

	return subscribersCopy
}

func (g *GoChannel) isClosed() bool {
	g.closedLock.Lock()
	defer g.closedLock.Unlock()

	return g.closed
}

// Close closes the GoChannel Pub/Sub.
func (g *GoChannel) Close() error {
	g.closedLock.Lock()
	defer g.closedLock.Unlock()

	if g.closed {
		return nil
	}

	g.closed = true
	close(g.closing)

	g.logger.Debug("Closing Pub/Sub, waiting for subscribers", nil)
	g.subscribersWg.Wait()

	g.logger.Info("Pub/Sub closed", nil)
	g.persistedMessages = nil

	return nil
}
