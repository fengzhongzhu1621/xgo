package metrics

import (
	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/prometheus/client_golang/prometheus"
)

var subscriberLabelKeys = []string{
	labelKeyHandlerName,
	labelKeySubscriberName,
}

// SubscriberPrometheusMetricsDecorator decorates a subscriber to capture Prometheus metrics.
type SubscriberPrometheusMetricsDecorator struct {
	message.Subscriber
	subscriberName                  string
	subscriberMessagesReceivedTotal *prometheus.CounterVec
	closing                         chan struct{}
}

func (s SubscriberPrometheusMetricsDecorator) recordMetrics(msg *message.Message) {
	if msg == nil {
		return
	}

	// 从消息上下文获取标签值
	ctx := msg.Context()
	labels := labelsFromCtx(ctx, subscriberLabelKeys...)
	if labels[labelKeySubscriberName] == "" {
		labels[labelKeySubscriberName] = s.subscriberName
	}
	if labels[labelKeyHandlerName] == "" {
		labels[labelKeyHandlerName] = labelValueNoHandler
	}

	go func() {
		if subscribeAlreadyObserved(ctx) {
			// decorator idempotency when applied decorator multiple times
			return
		}

		select {
		case <-msg.Acked():
			labels[labelAcked] = "acked"
		case <-msg.Nacked():
			labels[labelAcked] = "nacked"
		}
		// 上报指标记录订阅者已处理消息的总数
		s.subscriberMessagesReceivedTotal.With(labels).Inc()
	}()

	// 防止此decorator被重复加载
	msg.SetContext(setSubscribeObservedToCtx(msg.Context()))
}
