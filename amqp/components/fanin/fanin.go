package fanin

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/fengzhongzhu1621/xgo/log"
	"github.com/pkg/errors"
)

// FanIn is a component that receives messages from 1..N topics from a subscriber and publishes them
// on a specified topic in the publisher. In effect, messages are "multiplexed".
// 从router.Router派生的子类
type FanIn struct {
	router *router.Router
	config Config
	logger log.LoggerAdapter
}

// NewFanIn creates a new FanIn.
func NewFanIn(
	subscriber message.Subscriber,
	publisher message.Publisher,
	config Config,
	logger log.LoggerAdapter,
) (*FanIn, error) {
	if subscriber == nil {
		return nil, errors.New("missing subscriber")
	}
	if publisher == nil {
		return nil, errors.New("missing publisher")
	}

	config.SetDefaults()
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if logger == nil {
		logger = log.NopLogger{}
	}

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
	for _, topic := range config.SourceTopics {
		// 消息处理函数，将单个消息转换为消息数组
		handlerFunc := func(msg *message.Message) ([]*message.Message, error) {
			return []*message.Message{msg}, nil
		}
		router.AddHandler(
			fmt.Sprintf("fan_in_%s", topic), // handlerName
			topic,                           // subscribeTopic
			subscriber,
			config.TargetTopic, // publishTopic
			publisher,
			handlerFunc, // 消息处理函数
		)
	}

	return &FanIn{
		router: router,
		config: config,
		logger: logger,
	}, nil
}

// Run runs the FanIn.
func (f *FanIn) Run(ctx context.Context) error {
	return f.router.Run(ctx)
}

// Running is closed when FanIn is running.
func (f *FanIn) Running() chan struct{} {
	return f.router.Running()
}

// Close gracefully closes the FanIn
func (f *FanIn) Close() error {
	return f.router.Close()
}
