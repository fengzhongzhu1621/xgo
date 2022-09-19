package gochannel

import (
	"context"
	"sync"
	"xgo/log"
	"xgo/utils/randutils"

	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"

	"xgo/amqp/message"
)

// Config holds the GoChannel Pub/Sub's configuration options.
type Config struct {
	// Output channel buffer size.
	OutputChannelBuffer int64

	// If persistent is set to true, when subscriber subscribes to the topic,
	// it will receive all previously produced messages.
	//
	// All messages are persisted to the memory (simple slice),
	// so be aware that with large amount of messages you can go out of the memory.
	Persistent bool

	// When true, Publish will block until subscriber Ack's the message.
	// If there are no subscribers, Publish will not block (also when Persistent is true).
	BlockPublishUntilSubscriberAck bool
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
	logger log.LoggerAdapter

	subscribersWg          sync.WaitGroup
	subscribers            map[string][]*subscriber // 多个消费者
	subscribersLock        sync.RWMutex
	subscribersByTopicLock sync.Map // map of *sync.Mutex

	closed     bool
	closedLock sync.Mutex
	closing    chan struct{}

	persistedMessages     map[string][]*message.Message
	persistedMessagesLock sync.RWMutex
}

// NewGoChannel creates new GoChannel Pub/Sub.
//
// This GoChannel is not persistent.
// That means if you send a message to a topic to which no subscriber is subscribed, that message will be discarded.
func NewGoChannel(config Config, logger log.LoggerAdapter) *GoChannel {
	if logger == nil {
		logger = log.NopLogger{}
	}

	return &GoChannel{
		config: config,

		subscribers:            make(map[string][]*subscriber), // 多个订阅者
		subscribersByTopicLock: sync.Map{},
		logger: logger.With(log.LogFields{
			"pubsub_uuid": shortuuid.New(),
		}),

		closing: make(chan struct{}),

		persistedMessages: map[string][]*message.Message{},
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

	// 消息需要持久化存储
	if g.config.Persistent {
		g.persistedMessagesLock.Lock()
		if _, ok := g.persistedMessages[topic]; !ok {
			g.persistedMessages[topic] = make([]*message.Message, 0)
		}
		g.persistedMessages[topic] = append(g.persistedMessages[topic], messagesToPublish...)
		g.persistedMessagesLock.Unlock()
	}

	// 发送消息
	for i := range messagesToPublish {
		msg := messagesToPublish[i]

		ackedBySubscribers, err := g.sendMessage(topic, msg)
		if err != nil {
			return err
		}

		// 等待订阅者返回ack消息
		if g.config.BlockPublishUntilSubscriberAck {
			g.waitForAckFromSubscribers(msg, ackedBySubscribers)
		}
	}

	return nil
}

// waitForAckFromSubscribers 等待订阅者返回ack消息
func (g *GoChannel) waitForAckFromSubscribers(msg *message.Message, ackedByConsumer <-chan struct{}) {
	logFields := log.LogFields{"message_uuid": msg.UUID}
	g.logger.Debug("Waiting for subscribers ack", logFields)

	select {
	case <-ackedByConsumer:
		g.logger.Trace("Message acked by subscribers", logFields)
	case <-g.closing:
		g.logger.Trace("Closing Pub/Sub before ack from subscribers", logFields)
	}
}

func (g *GoChannel) sendMessage(topic string, message *message.Message) (<-chan struct{}, error) {
	// 选择多个订阅者
	subscribers := g.topicSubscribers(topic)
	ackedBySubscribers := make(chan struct{})

	logFields := log.LogFields{"message_uuid": message.UUID, "topic": topic}

	if len(subscribers) == 0 {
		close(ackedBySubscribers)
		g.logger.Info("No subscribers to send message", logFields)
		return ackedBySubscribers, nil
	}

	go func(subscribers []*subscriber) {
		wg := &sync.WaitGroup{}

		for i := range subscribers {
			subscriber := subscribers[i]

			wg.Add(1)
			go func() {
				// 每个订阅者都收到消息
				subscriber.sendMessageToSubscriber(message, logFields)
				wg.Done()
			}()
		}

		wg.Wait()
		close(ackedBySubscribers)
	}(subscribers)

	return ackedBySubscribers, nil
}

// Subscribe returns channel to which all published messages are sent.
// Messages are not persisted. If there are no subscribers and message is produced it will be gone.
//
// There are no consumer groups support etc. Every consumer will receive every produced message.
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

	s := &subscriber{
		ctx:           ctx,
		uuid:          randutils.NewUUID(),
		outputChannel: make(chan *message.Message, g.config.OutputChannelBuffer),
		logger:        g.logger,
		closing:       make(chan struct{}),
	}

	go func(s *subscriber, g *GoChannel) {
		select {
		case <-ctx.Done():
			// unblock
		case <-g.closing:
			// unblock
		}

		s.Close()

		g.subscribersLock.Lock()
		defer g.subscribersLock.Unlock()

		subLock, _ := g.subscribersByTopicLock.Load(topic)
		subLock.(*sync.Mutex).Lock()
		defer subLock.(*sync.Mutex).Unlock()

		g.removeSubscriber(topic, s)
		g.subscribersWg.Done()
	}(s, g)

	if !g.config.Persistent {
		defer g.subscribersLock.Unlock()
		defer subLock.(*sync.Mutex).Unlock()

		g.addSubscriber(topic, s)

		return s.outputChannel, nil
	}

	go func(s *subscriber) {
		defer g.subscribersLock.Unlock()
		defer subLock.(*sync.Mutex).Unlock()

		g.persistedMessagesLock.RLock()
		messages, ok := g.persistedMessages[topic]
		g.persistedMessagesLock.RUnlock()

		if ok {
			for i := range messages {
				msg := g.persistedMessages[topic][i]
				logFields := log.LogFields{"message_uuid": msg.UUID, "topic": topic}

				go s.sendMessageToSubscriber(msg, logFields)
			}
		}

		g.addSubscriber(topic, s)
	}(s)

	return s.outputChannel, nil
}

func (g *GoChannel) addSubscriber(topic string, s *subscriber) {
	if _, ok := g.subscribers[topic]; !ok {
		g.subscribers[topic] = make([]*subscriber, 0)
	}
	g.subscribers[topic] = append(g.subscribers[topic], s)
}

func (g *GoChannel) removeSubscriber(topic string, toRemove *subscriber) {
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
func (g *GoChannel) topicSubscribers(topic string) []*subscriber {
	subscribers, ok := g.subscribers[topic]
	if !ok {
		return nil
	}

	// let's do a copy to avoid race conditions and deadlocks due to lock
	// 复制订阅者
	subscribersCopy := make([]*subscriber, len(subscribers))
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
