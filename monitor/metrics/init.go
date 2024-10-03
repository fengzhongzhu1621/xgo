package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	callerInitMetricsOnce sync.Once
	serviceName           = "xgo"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "requests_total",
			Help:        "requests_total",
			ConstLabels: prometheus.Labels{"service": serviceName},
		},
		[]string{"method", "path", "status", "error", "client_id"},
	)

	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "request_duration_milliseconds",
		Help:        "request_duration_milliseconds",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     []float64{50, 100, 200, 500, 1000, 2000, 5000},
	},
		[]string{"method", "path", "status", "client_id"},
	)
)

func InitMetrics() {
	callerInitMetricsOnce.Do(func() {
		prometheus.MustRegister(RequestCount)
		prometheus.MustRegister(RequestDuration)
	})
}
