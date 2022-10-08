package main

import (
	"context"
	"errors"
	"flag"
	"math"
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var (
	metricsAddr  = flag.String("metrics", ":8081", "The address that will expose /metrics for Prometheus")
	handlerDelay = flag.Float64("delay", 0, "The stdev of normal distribution of delay in handler (in seconds), to simulate load")

	logger = watermill.NewStdLogger(true, true)
	// 设置随机种子
	random = rand.New(rand.NewSource(time.Now().Unix()))
)

// 延迟n秒，符合正态分布
func delay() {
	seconds := *handlerDelay
	if seconds == 0 {
		return
	}
	// NormFloat64: 返回一个服从标准正态分布（标准差=1，期望=0）、
	// 取值范围在[-math.MaxFloat64, +math.MaxFloat64]的float64值。
	// 如果要生成不同的正态分布值，调用者可用如下代码调整输出
	// sample = NormFloat64() * 标准差 + 期望
	delay := math.Abs(random.NormFloat64() * seconds)
	time.Sleep(time.Duration(float64(time.Second) * delay))
}

// handler publishes 0-4 messages with a random delay.
func handler(msg *message.Message) ([]*message.Message, error) {
	delay()

	numOutgoing := random.Intn(4)
	outgoing := make([]*message.Message, numOutgoing)

	for i := 0; i < numOutgoing; i++ {
		outgoing[i] = msg.Copy()
	}
	return outgoing, nil
}

// consumeMessages consumes the messages exiting the handler.
func consumeMessages(subscriber message.Subscriber) {
	messages, err := subscriber.Subscribe(context.Background(), "pub_topic")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		msg.Ack()
	}
}

// produceMessages produces the incoming messages in delays of 50-100 milliseconds.
func produceMessages(routerClosed chan struct{}, publisher message.Publisher) {
	for {
		select {
		case <-routerClosed:
			return
		default:
			// go on
		}

		time.Sleep(50*time.Millisecond + time.Duration(random.Intn(50))*time.Millisecond)
		msg := message.NewMessage(watermill.NewUUID(), []byte{})
		_ = publisher.Publish("sub_topic", msg)
	}
}

func main() {
	flag.Parse()

	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router, err := message.NewRouter(
		message.RouterConfig{},
		logger,
	)
	if err != nil {
		panic(err)
	}

	// 创建metrics httpsever
	prometheusRegistry, closeMetricsServer := metrics.CreateRegistryAndServeHTTP(*metricsAddr)
	defer closeMetricsServer()

	// 创建metric builder对象
	// we leave the namespace and subsystem empty
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(prometheusRegistry, "", "")
	metricsBuilder.AddPrometheusRouterMetrics(router)

	router.AddMiddleware(
		middleware.Recoverer,
		middleware.RandomFail(0.1),
		middleware.RandomPanic(0.1),
	)
	router.AddPlugin(plugin.SignalsHandler)

	router.AddHandler(
		"metrics-example",
		"sub_topic",
		pubSub,
		"pub_topic",
		pubSub,
		handler,
	)

	pub := randomFailPublisherDecorator{pubSub, 0.1}

	// The handler's publisher and subscriber will be decorated by `AddPrometheusRouterMetrics`.
	// We are using the same pub/sub to generate messages incoming to the handler
	// and consume the outgoing messages.
	// They will have `handler_name=<no handler>` label in Prometheus.
	// 生成一个新的生产者，指向同一个gochannel对象
	subWithMetrics, err := metricsBuilder.DecorateSubscriber(pubSub)
	if err != nil {
		panic(err)
	}
	// 生成一个新的消费者，指向同一个gochannel对象
	pubWithMetrics, err := metricsBuilder.DecoratePublisher(pub)
	if err != nil {
		panic(err)
	}

	routerClosed := make(chan struct{})
	go produceMessages(routerClosed, pubWithMetrics)
	go consumeMessages(subWithMetrics)

	_ = router.Run(context.Background())
	close(routerClosed)
}

type randomFailPublisherDecorator struct {
	message.Publisher
	failProbability float64
}

func (r randomFailPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	if random.Float64() < r.failProbability {
		return errors.New("random publishing failure")
	}
	return r.Publisher.Publish(topic, messages...)
}
