package message

import (
	"context"
	"sync"
)

// PublisherDecorator wraps the underlying Publisher, adding some functionality.
type PublisherDecorator func(pub Publisher) (Publisher, error)

// SubscriberDecorator wraps the underlying Subscriber, adding some functionality.
type SubscriberDecorator func(sub Subscriber) (Subscriber, error)

type messageTransformSubscriberDecorator struct {
	sub Subscriber // 订阅者

	transform   func(*Message) // 转换函数
	subscribeWg sync.WaitGroup
}

// MessageTransformSubscriberDecorator creates a subscriber decorator that calls transform
// on each message that passes through the subscriber.
// 方便每个订阅者都要处理同一个消息
func MessageTransformSubscriberDecorator(transform func(*Message)) SubscriberDecorator {
	if transform == nil {
		panic("transform function is nil")
	}
	// 每个Subscriber处理同一个消息
	return func(sub Subscriber) (Subscriber, error) {
		return &messageTransformSubscriberDecorator{
			sub:       sub,
			transform: transform,
		}, nil
	}
}

func (t *messageTransformSubscriberDecorator) Subscribe(ctx context.Context, topic string) (<-chan *Message, error) {
	// 订阅topic
	in, err := t.sub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	out := make(chan *Message)
	t.subscribeWg.Add(1)
	go func() {
		// 从管道消费消息
		for msg := range in {
			// 格式化消息
			t.transform(msg)
			out <- msg
		}
		// 关闭订阅者后执行
		close(out)
		t.subscribeWg.Done()
	}()

	return out, nil
}

func (t *messageTransformSubscriberDecorator) Close() error {
	// 关闭订阅者
	err := t.sub.Close()
	// 等待正在发送的订阅者结束消息发送
	t.subscribeWg.Wait()
	return err
}

type messageTransformPublisherDecorator struct {
	Publisher
	transform func(*Message)
}

// MessageTransformPublisherDecorator creates a publisher decorator that calls transform
// on each message that passes through the publisher.
func MessageTransformPublisherDecorator(transform func(*Message)) PublisherDecorator {
	if transform == nil {
		panic("transform function is nil")
	}
	return func(pub Publisher) (Publisher, error) {
		return &messageTransformPublisherDecorator{
			Publisher: pub,
			transform: transform,
		}, nil
	}
}

// Publish applies the transform to each message and returns the underlying Publisher's result.
func (d messageTransformPublisherDecorator) Publish(topic string, messages ...*Message) error {
	for i := range messages {
		// 格式化消息
		d.transform(messages[i])
	}
	// 发送消息到队列
	return d.Publisher.Publish(topic, messages...)
}
