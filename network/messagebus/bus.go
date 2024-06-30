package messagebus

import (
	"fmt"
	"reflect"
	"sync"
)

// MessageBus implements publish/subscribe messaging paradigm
type MessageBus interface {
	// Publish publishes arguments to the given topic subscribers
	// Publish block only when the buffer of one of the subscribers is full.
	Publish(topic string, args ...interface{})
	// Close unsubscribe all handlers from given topic
	Close(topic string)
	// Subscribe subscribes to the given topic
	Subscribe(topic string, fn interface{}) error
	// Unsubscribe unsubscribe handler from the given topic
	Unsubscribe(topic string, fn interface{}) error
}

// 存放所有的订阅者
type handlersMap map[string][]*handler

// 订阅者
type handler struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

type messageBus struct {
	handlerQueueSize int // 订阅者消息队列的大小
	mtx              sync.RWMutex
	handlers         handlersMap
}

func (b *messageBus) Publish(topic string, args ...interface{}) {
	// 构造消费者函数的参数，转换为reflect.Value格式
	rArgs := buildHandlerArgs(args)

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	// 根据topic找到所有的订阅者
	if hs, ok := b.handlers[topic]; ok {
		for _, h := range hs {
			// 将函数的参数添加到消费队列
			h.queue <- rArgs
		}
	}
}

// 订阅消息topic
func (b *messageBus) Subscribe(topic string, fn interface{}) error {
	// fn 必须是函数
	if err := isValidHandler(fn); err != nil {
		return err
	}

	// 创建一个订阅者
	h := &handler{
		callback: reflect.ValueOf(fn),                            // 消费者函数
		queue:    make(chan []reflect.Value, b.handlerQueueSize), // 每个订阅者有自己的函数参数队列
	}

	// 注册订阅者消费消息协程
	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	// 注册订阅者
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)

	return nil
}

func (b *messageBus) Unsubscribe(topic string, fn interface{}) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	rv := reflect.ValueOf(fn)

	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for i, h := range b.handlers[topic] {
			if h.callback == rv {
				// 关闭某主题下所有的订阅者参数通道
				close(h.queue)
				// 删除该订阅者
				if len(b.handlers[topic]) == 1 {
					delete(b.handlers, topic)
				} else {
					b.handlers[topic] = append(b.handlers[topic][:i], b.handlers[topic][i+1:]...)
				}
			}
		}

		return nil
	}

	return fmt.Errorf("topic %s doesn't exist", topic)
}

func (b *messageBus) Close(topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for _, h := range b.handlers[topic] {
			close(h.queue)
		}

		delete(b.handlers, topic)

		return
	}
}

// 订阅处理必须是函数
func isValidHandler(fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return nil
}

// 构造消费者函数的参数
func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}

// New creates new MessageBus
// handlerQueueSize sets buffered channel length per subscriber
func New(handlerQueueSize int) MessageBus {
	if handlerQueueSize == 0 {
		panic("handlerQueueSize has to be greater then 0")
	}

	return &messageBus{
		handlerQueueSize: handlerQueueSize,
		handlers:         make(handlersMap),
	}
}
