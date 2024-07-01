package metrics

import (
	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/fengzhongzhu1621/xgo/buildin"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPrometheusMetricsBuilder(prometheusRegistry prometheus.Registerer, namespace string, subsystem string) PrometheusMetricsBuilder {
	return PrometheusMetricsBuilder{
		Namespace:          namespace,
		Subsystem:          subsystem,
		PrometheusRegistry: prometheusRegistry,
	}
}

// PrometheusMetricsBuilder provides methods to decorate publishers, subscribers and handlers.
type PrometheusMetricsBuilder struct {
	// PrometheusRegistry may be filled with a pre-existing Prometheus registry, or left empty for the default registry.
	PrometheusRegistry prometheus.Registerer

	Namespace string
	Subsystem string
}

// AddPrometheusRouterMetrics is a convenience function that acts on the message router to add the metrics middleware
// to all its handlers. The handlers' publishers and subscribers are also decorated.
func (b PrometheusMetricsBuilder) AddPrometheusRouterMetrics(r *router.Router) {
	r.AddPublisherDecorators(b.DecoratePublisher)
	// 注册 counter metrics，记录已处理消息的总数
	r.AddSubscriberDecorators(b.DecorateSubscriber)
	// 注册 Histogram metrics
	r.AddMiddleware(b.NewRouterMiddleware().Middleware)
}

// DecoratePublisher wraps the underlying publisher with Prometheus metrics.
func (b PrometheusMetricsBuilder) DecoratePublisher(pub message.Publisher) (message.Publisher, error) {
	var err error
	d := PublisherPrometheusMetricsDecorator{
		pub:           pub,
		publisherName: buildin.StructName(pub),
	}

	d.publishTimeSeconds, err = b.registerHistogramVec(prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: b.Namespace,
			Subsystem: b.Subsystem,
			Name:      "publish_time_seconds",
			Help:      "The time that a publishing attempt (success or not) took in seconds",
		},
		publisherLabelKeys,
	))
	if err != nil {
		return nil, errors.Wrap(err, "could not register publish time metric")
	}
	return d, nil
}

// DecorateSubscriber wraps the underlying subscriber with Prometheus metrics.
func (b PrometheusMetricsBuilder) DecorateSubscriber(sub message.Subscriber) (message.Subscriber, error) {
	var err error
	d := &SubscriberPrometheusMetricsDecorator{
		closing:        make(chan struct{}),
		subscriberName: buildin.StructName(sub),
	}

	d.subscriberMessagesReceivedTotal, err = b.registerCounterVec(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: b.Namespace,
			Subsystem: b.Subsystem,
			Name:      "subscriber_messages_received_total",
			Help:      "The total number of messages received by the subscriber",
		},
		append(subscriberLabelKeys, labelAcked),
	))
	if err != nil {
		return nil, errors.Wrap(err, "could not register time to ack metric")
	}

	// 给订阅者添加 transform 方法
	d.Subscriber, err = message.MessageTransformSubscriberDecorator(d.recordMetrics)(sub)
	if err != nil {
		return nil, errors.Wrap(err, "could not decorate subscriber with metrics decorator")
	}

	return d, nil
}

func (b PrometheusMetricsBuilder) register(c prometheus.Collector) (prometheus.Collector, error) {
	err := b.PrometheusRegistry.Register(c)
	if err == nil {
		return c, nil
	}

	if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
		return are.ExistingCollector, nil
	}

	return nil, err
}

func (b PrometheusMetricsBuilder) registerCounterVec(c *prometheus.CounterVec) (*prometheus.CounterVec, error) {
	col, err := b.register(c)
	if err != nil {
		return nil, err
	}
	return col.(*prometheus.CounterVec), nil
}

func (b PrometheusMetricsBuilder) registerHistogramVec(h *prometheus.HistogramVec) (*prometheus.HistogramVec, error) {
	col, err := b.register(h)
	if err != nil {
		return nil, err
	}
	return col.(*prometheus.HistogramVec), nil
}

// NewRouterMiddleware returns new middleware.
// handler执行前先注册metrics
func (b PrometheusMetricsBuilder) NewRouterMiddleware() HandlerPrometheusMetricsMiddleware {
	var err error
	m := HandlerPrometheusMetricsMiddleware{}

	m.handlerExecutionTimeSeconds, err = b.registerHistogramVec(prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: b.Namespace,
			Subsystem: b.Subsystem,
			Name:      "handler_execution_time_seconds",
			Help:      "The total time elapsed while executing the handler function in seconds",
			Buckets:   handlerExecutionTimeBuckets,
		},
		handlerLabelKeys,
	))
	if err != nil {
		panic(errors.Wrap(err, "could not register handler execution time metric"))
	}

	return m
}
