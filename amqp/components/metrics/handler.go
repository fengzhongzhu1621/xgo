package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fengzhongzhu1621/xgo/amqp/message"
	"github.com/fengzhongzhu1621/xgo/amqp/router"
)

var (
	handlerLabelKeys = []string{
		labelKeyHandlerName,
		labelSuccess,
	}

	// handlerExecutionTimeBuckets are one order of magnitude smaller than default buckets (5ms~10s),
	// because the handler execution times are typically shorter (µs~ms range).
	handlerExecutionTimeBuckets = []float64{
		0.0005,
		0.001,
		0.0025,
		0.005,
		0.01,
		0.025,
		0.05,
		0.1,
		0.25,
		0.5,
		1,
	}
)

// HandlerPrometheusMetricsMiddleware is a middleware that captures Prometheus metrics.
type HandlerPrometheusMetricsMiddleware struct {
	handlerExecutionTimeSeconds *prometheus.HistogramVec
}

// Middleware returns the middleware ready to be used with watermill's Router.
func (m HandlerPrometheusMetricsMiddleware) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		now := time.Now()
		ctx := msg.Context()
		labels := prometheus.Labels{
			labelKeyHandlerName: router.HandlerNameFromCtx(ctx),
		}

		defer func() {
			if err != nil {
				labels[labelSuccess] = "false"
			} else {
				labels[labelSuccess] = "true"
			}
			// 消息处理的耗时（包括中间件）
			m.handlerExecutionTimeSeconds.With(labels).Observe(time.Since(now).Seconds())
		}()

		return h(msg)
	}
}
