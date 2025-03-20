package rest

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/collections/flowctrl"
	"github.com/fengzhongzhu1621/xgo/ginx/discovery"
	"github.com/fengzhongzhu1621/xgo/ginx/types"
	"github.com/prometheus/client_golang/prometheus"
)

type Capability struct {
	Client     IHttpClient
	Discover   discovery.Interface
	Throttle   flowctrl.RateLimiter
	Mock       types.MockInfo
	MetricOpts MetricOption
	// the max tolerance api request latency time, if exceeded this time, then
	// this request will be logged and warned.
	ToleranceLatencyTime time.Duration
}

type MetricOption struct {
	// prometheus metric register
	Register prometheus.Registerer
	// if not set, use default buckets value
	DurationBuckets []float64
}
