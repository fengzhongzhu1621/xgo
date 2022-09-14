package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	. "xgo/amqp/message"
	. "xgo/amqp/publisher"
	"xgo/log"
	"xgo/utils"
	sync_internal "xgo/utils/sync"
)

// NoPublishHandlerFunc is HandlerFunc alternative, which doesn't produce any messages.
type NoPublishHandlerFunc func(msg *Message) error

// PassthroughHandler is a handler that passes the message unchanged from the subscriber to the publisher.
var PassthroughHandler HandlerFunc = func(msg *Message) ([]*Message, error) {
	return []*Message{msg}, nil
}

// RouterPlugin is function which is executed on Router start.
type RouterPlugin func(*Router) error

// Router is responsible for handling messages from subscribers using provided handler functions.
//
// If the handler function returns a message, the message is published with the publisher.
// You can use middlewares to wrap handlers with common logic like logging, instrumentation, etc.
type Router struct {
	config RouterConfig // 路由配置

	middlewares []Middleware // 包含多个消息中间件

	plugins []RouterPlugin // 路由的前置逻辑

	handlers     map[string]*handler // 路由处理器
	handlersLock sync.RWMutex        // 路由处理器新增、关闭和执行时加锁

	handlersWg        *sync.WaitGroup // 关闭路由时判断所有的handler执行完毕
	runningHandlersWg *sync.WaitGroup // 标记所有handler接收到的一批消息是否都处理完毕
	handlerAdded      chan struct{}   // 当添加一个handler时AddHandler()触发此事件

	closeCh    chan struct{} // 路由接收到的关闭信号
	closedCh   chan struct{} // 路由接收到的超时关闭信号
	closed     bool          // 路由是否已经关闭
	closedLock sync.Mutex

	logger log.LoggerAdapter

	publisherDecorators  []PublisherDecorator  // 生产者装饰器
	subscriberDecorators []SubscriberDecorator // 消费者装饰器

	isRunning bool          // 用于判断router是否正在运行
	running   chan struct{} // 路由执行时告诉所有的处理器协程开始执行
}

// NewRouter creates a new Router with given configuration.
func NewRouter(config RouterConfig, logger log.LoggerAdapter) (*Router, error) {
	// 路由配置初始化默认值
	config.setDefaults()
	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config")
	}

	return &Router{
		config: config, // 路由配置

		handlers: map[string]*handler{}, // 包含多个处理器

		handlersWg:        &sync.WaitGroup{},
		runningHandlersWg: &sync.WaitGroup{},
		handlerAdded:      make(chan struct{}),

		closeCh:  make(chan struct{}), // 路由接收到的关闭信号
		closedCh: make(chan struct{}), // 路由接收到的超时关闭信号

		logger: logger, // 日志处理器

		running: make(chan struct{}), // 路由执行时告诉所有的处理器协程开始执行
	}, nil
}

// Logger returns the Router's logger.
func (r *Router) Logger() log.LoggerAdapter {
	return r.logger
}

// AddMiddleware adds a new middleware to the router.
//
// The order of middleware matters. Middleware added at the beginning is executed first.
func (r *Router) AddMiddleware(m ...HandlerMiddleware) {
	r.logger.Debug("Adding middleware", log.LogFields{"count": fmt.Sprintf("%d", len(m))})

	r.addRouterLevelMiddleware(m...)
}

// 添加路由级别的中间件
func (r *Router) addRouterLevelMiddleware(m ...HandlerMiddleware) {
	for _, handlerMiddleware := range m {
		middleware := Middleware{
			Handler:       handlerMiddleware,
			HandlerName:   "", // 路由级别的消息中间件没有名字
			IsRouterLevel: true,
		}
		r.middlewares = append(r.middlewares, middleware)
	}
}

// 添加handler级别的中间件
func (r *Router) addHandlerLevelMiddleware(handlerName string, m ...HandlerMiddleware) {
	for _, handlerMiddleware := range m {
		middleware := Middleware{
			Handler:       handlerMiddleware,
			HandlerName:   handlerName,
			IsRouterLevel: false,
		}
		r.middlewares = append(r.middlewares, middleware)
	}
}

// AddPlugin adds a new plugin to the router.
// Plugins are executed during startup of the router.
//
// A plugin can, for example, close the router after SIGINT or SIGTERM is sent to the process (SignalsHandler plugin).
func (r *Router) AddPlugin(p ...RouterPlugin) {
	r.logger.Debug("Adding plugins", log.LogFields{"count": fmt.Sprintf("%d", len(p))})

	r.plugins = append(r.plugins, p...)
}

// AddPublisherDecorators wraps the router's Publisher.
// The first decorator is the innermost, i.e. calls the original publisher.
// 添加生产者装饰器
func (r *Router) AddPublisherDecorators(dec ...PublisherDecorator) {
	r.logger.Debug("Adding publisher decorators", log.LogFields{"count": fmt.Sprintf("%d", len(dec))})

	r.publisherDecorators = append(r.publisherDecorators, dec...)
}

// AddSubscriberDecorators wraps the router's Subscriber.
// The first decorator is the innermost, i.e. calls the original subscriber.
// 添加消费者装饰器
func (r *Router) AddSubscriberDecorators(dec ...SubscriberDecorator) {
	r.logger.Debug("Adding subscriber decorators", log.LogFields{"count": fmt.Sprintf("%d", len(dec))})

	r.subscriberDecorators = append(r.subscriberDecorators, dec...)
}

// Handlers returns all registered handlers.
// 返回所有已注册的处理器函数
func (r *Router) Handlers() map[string]HandlerFunc {
	handlers := map[string]HandlerFunc{}

	for handlerName, handler := range r.handlers {
		handlers[handlerName] = handler.handlerFunc
	}

	return handlers
}

// AddHandler adds a new handler.
//
// handlerName must be unique. For now, it is used only for debugging.
//
// subscribeTopic is a topic from which handler will receive messages.
//
// publishTopic is a topic to which router will produce messages returned by handlerFunc.
// When handler needs to publish to multiple topics,
// it is recommended to just inject Publisher to Handler or implement middleware
// which will catch messages and publish to topic based on metadata for example.
//
// If handler is added while router is already running, you need to explicitly call RunHandlers().
func (r *Router) AddHandler(
	handlerName string, // 消息处理器名称
	subscribeTopic string, // 消费者topic
	subscriber Subscriber, // 消费者
	publishTopic string, // 生产者topic
	publisher Publisher, // 生产者，消费后的数据重新发给指定生产者
	handlerFunc HandlerFunc, // 消息处理函数
) *Handler {
	r.logger.Info("Adding handler", log.LogFields{
		"handler_name": handlerName,
		"topic":        subscribeTopic,
	})

	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()

	// 消息处理器不能重复添加
	if _, ok := r.handlers[handlerName]; ok {
		panic(DuplicateHandlerNameError{handlerName})
	}
	// 获得生产者和消费者的结构体名称
	publisherName, subscriberName := utils.StructName(publisher), utils.StructName(subscriber)

	newHandler := &handler{
		name:   handlerName,
		logger: r.logger,

		subscriber:     subscriber,     // 消费者
		subscribeTopic: subscribeTopic, // 消费者topic
		subscriberName: subscriberName, // 消费者名称

		publisher:     publisher,
		publishTopic:  publishTopic,
		publisherName: publisherName,

		handlerFunc:       handlerFunc,         // 消息处理函数
		runningHandlersWg: r.runningHandlersWg, // 标记所有handler接收到的一批消息是否都处理完毕
		messagesCh:        nil,
		routersCloseCh:    r.closeCh, // 路由关闭信号

		startedCh: make(chan struct{}), // handler开始执行信号
	}

	// 用于等待一个消息关联的所有的handlers执行完成
	r.handlersWg.Add(1)

	r.handlers[handlerName] = newHandler

	// 触发一个处理器添加事件
	select {
	case r.handlerAdded <- struct{}{}:
	default:
		// closeWhenAllHandlersStopped is not always waiting for handlerAdded
	}

	// 返回路由和处理器之间的关系
	return &Handler{
		router:  r,
		handler: newHandler,
	}
}

// AddNoPublisherHandler adds a new handler.
// This handler cannot return messages.
// When message is returned it will occur an error and Nack will be sent.
//
// handlerName must be unique. For now, it is used only for debugging.
//
// subscribeTopic is a topic from which handler will receive messages.
//
// subscriber is Subscriber from which messages will be consumed.
//
// If handler is added while router is already running, you need to explicitly call RunHandlers().
func (r *Router) AddNoPublisherHandler(
	handlerName string,
	subscribeTopic string,
	subscriber Subscriber,
	handlerFunc NoPublishHandlerFunc,
) *Handler {
	// 消息处理函数
	handlerFuncAdapter := func(msg *Message) ([]*Message, error) {
		return nil, handlerFunc(msg)
	}

	return r.AddHandler(handlerName, subscribeTopic, subscriber, "", DisabledPublisher{}, handlerFuncAdapter)
}

// Run runs all plugins and handlers and starts subscribing to provided topics.
// This call is blocking while the router is running.
//
// When all handlers have stopped (for example, because subscriptions were closed), the router will also stop.
//
// To stop Run() you should call Close() on the router.
//
// ctx will be propagated to all subscribers.
//
// When all handlers are stopped (for example: because of closed connection), Run() will be also stopped.
func (r *Router) Run(ctx context.Context) (err error) {
	// 标记router正在执行，防止重复执行
	if r.isRunning {
		return errors.New("router is already running")
	}
	r.isRunning = true

	// 1. WithCancel()函数接受一个 Context 并返回其子Context和取消函数cancel
	// 2. 新创建协程中传入子Context做参数，且需监控子Context的Done通道，若收到消息，则退出
	// 3. 需要新协程结束时，在外面调用 cancel 函数，即会往子Context的Done通道发送消息
	// 4. 注意：当 父Context的 Done() 关闭的时候，子 ctx 的 Done() 也会被关闭
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 执行处理器前置操作
	r.logger.Debug("Loading plugins", nil)
	for _, plugin := range r.plugins {
		if err := plugin(r); err != nil {
			return errors.Wrapf(err, "cannot initialize plugin %v", plugin)
		}
	}

	// 创建处理器协程
	if err := r.RunHandlers(ctx); err != nil {
		return err
	}

	// 标记路由正在执行 和 参数isRunning的功能类似
	// 不同之处在于可以告诉其他协程
	close(r.running)

	// 执行自动关闭协程，等待所有的处理器执行完成并关闭路由，防止忘记调用Close()
	// 使用协程的原因是，关闭事件closeCh可以由其它地方触发Close()方法执行关闭
	go r.closeWhenAllHandlersStopped()

	// 等待路由关闭事件，由Close()触发
	<-r.closeCh

	// 给处理器发送上下文关闭信号
	cancel()

	r.logger.Info("Waiting for messages", log.LogFields{
		"timeout": r.config.CloseTimeout,
	})

	// 等待路由超时关闭信号
	<-r.closedCh

	r.logger.Info("All messages processed", nil)

	return nil
}

// RunHandlers runs all handlers that were added after Run().
// RunHandlers is idempotent, so can be called multiple times safely.
func (r *Router) RunHandlers(ctx context.Context) error {
	if !r.isRunning {
		return errors.New("you can't call RunHandlers on non-running router")
	}

	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()

	// 遍历所有的处理器对象
	for name, h := range r.handlers {
		name := name
		h := h

		// 判断处理器是否已经执行
		if h.started {
			continue
		}

		// 给handler的生产者逆序添加装饰器
		if err := r.decorateHandlerPublisher(h); err != nil {
			return errors.Wrapf(err, "could not decorate publisher of handler %s", name)
		}
		// 给handler的消费者顺序添加装饰器
		if err := r.decorateHandlerSubscriber(h); err != nil {
			return errors.Wrapf(err, "could not decorate subscriber of handler %s", name)
		}

		r.logger.Debug("Subscribing to topic", log.LogFields{
			"subscriber_name": h.name,
			"topic":           h.subscribeTopic,
		})

		ctx, cancel := context.WithCancel(ctx)

		// 从队列的指定topic获取消息，通常来说，不同的handler处理不同的topic
		messages, err := h.subscriber.Subscribe(ctx, h.subscribeTopic)
		if err != nil {
			cancel()
			return errors.Wrapf(err, "cannot subscribe topic %s", h.subscribeTopic)
		}

		// 从队列读取消息，将格式化后的的消息放入此队列
		h.messagesCh = messages
		// 标记handler开始执行
		h.started = true
		// 发送消息处理器开始执行信号
		close(h.startedCh)

		h.stopFn = cancel
		h.stopped = make(chan struct{})

		// 创建handler协程，当router发送 close(r.running) 信号后协程开始执行
		go func() {
			// handler协程执行完毕后，给subscriber发送关闭信号
			defer cancel()

			// 执行handler
			h.run(ctx, r.middlewares)

			// 告诉router，handler执行完毕
			r.handlersWg.Done()
			r.logger.Info("Subscriber stopped", log.LogFields{
				"subscriber_name": h.name,
				"topic":           h.subscribeTopic,
			})

			r.handlersLock.Lock()
			delete(r.handlers, name)
			r.handlersLock.Unlock()
		}()
	}
	return nil
}

// closeWhenAllHandlersStopped closed router, when all handlers has stopped,
// because for example all subscriptions are closed.
func (r *Router) closeWhenAllHandlersStopped() {
	r.handlersLock.RLock()
	noHandlers := len(r.handlers) == 0
	r.handlersLock.RUnlock()

	if noHandlers {
		// in that situation router will be closed immediately (even if they are no routers)
		// let's wait for
		select {
		case <-r.handlerAdded:
			// it should be some handler to track
		case <-r.closedCh:
			// let's avoid goroutine leak
			return
		}
	}
	// 等待所有的handers执行完成
	r.handlersWg.Wait()
	if r.isClosed() {
		// already closed
		return
	}

	r.logger.Error("All handlers stopped, closing router", errors.New("all router handlers stopped"), nil)

	if err := r.Close(); err != nil {
		r.logger.Error("Cannot close router", err, nil)
	}
}

// Running is closed when router is running.
// In other words: you can wait till router is running using
//		fmt.Println("Starting router")
//		go r.Run(ctx)
//		<- r.Running()
//		fmt.Println("Router is running")
func (r *Router) Running() chan struct{} {
	return r.running
}

// IsRunning returns true when router is running.
func (r *Router) IsRunning() bool {
	select {
	case <-r.running:
		// 路由正在执行
		return true
	default:
		// 路由还没有执行
		return false
	}
}

// Close gracefully closes the router with a timeout provided in the configuration.
func (r *Router) Close() error {
	r.closedLock.Lock()
	defer r.closedLock.Unlock()

	r.handlersLock.Lock()
	defer r.handlersLock.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true

	r.logger.Info("Closing router", nil)
	defer r.logger.Info("Router closed", nil)

	close(r.closeCh)
	defer close(r.closedCh)

	timeouted := sync_internal.WaitGroupTimeout(r.handlersWg, r.config.CloseTimeout)
	if timeouted {
		return errors.New("router close timeout")
	}

	return nil
}

func (r *Router) isClosed() bool {
	r.closedLock.Lock()
	defer r.closedLock.Unlock()

	return r.closed
}

// decorateHandlerPublisher applies the decorator chain to handler's publisher.
// They are applied in reverse order, so that the later decorators use the result of former ones.
func (r *Router) decorateHandlerPublisher(h *handler) error {
	var err error
	pub := h.publisher
	for i := len(r.publisherDecorators) - 1; i >= 0; i-- {
		decorator := r.publisherDecorators[i]
		pub, err = decorator(pub)
		if err != nil {
			return errors.Wrap(err, "could not apply publisher decorator")
		}
	}
	r.handlers[h.name].publisher = pub
	return nil
}

// decorateHandlerSubscriber applies the decorator chain to handler's subscriber.
// They are applied in regular order, so that the later decorators use the result of former ones.
func (r *Router) decorateHandlerSubscriber(h *handler) error {
	var err error
	sub := h.subscriber

	// add values to message context to subscriber
	// it goes before other decorators, so that they may take advantage of these values
	messageTransform := func(msg *Message) {
		if msg != nil {
			h.addHandlerContext(msg)
		}
	}
	// 给消费者sub添加transform方法
	sub, err = MessageTransformSubscriberDecorator(messageTransform)(sub)
	if err != nil {
		return errors.Wrapf(err, "cannot wrap subscriber with context decorator")
	}

	for _, decorator := range r.subscriberDecorators {
		sub, err = decorator(sub)
		if err != nil {
			return errors.Wrap(err, "could not apply subscriber decorator")
		}
	}
	r.handlers[h.name].subscriber = sub
	return nil
}
