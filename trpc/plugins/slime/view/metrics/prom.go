package metrics

import (
	"github.com/fengzhongzhu1621/xgo/opentelemetry/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// 子系统前缀：subSystem = "slime" 确保指标命名唯一性（如 slime_appRequest_total）
	subSystem = "slime"
)

// Emitter defines a common metric interface.
type IEmitter interface {
	metrics.ICounter
	metrics.IHistogram
}

// Emitter is a prometheus Reporter.
type Emitter struct {
	appReq   *prometheus.CounterVec   // 应用层请求计数器（如用户请求数）
	realReq  *prometheus.CounterVec   // 真实请求计数器（如RPC调用数）
	appCost  *prometheus.HistogramVec // 应用层耗时直方图（端到端延迟）
	realCost *prometheus.HistogramVec // 真实请求耗时直方图（单次调用延迟）
}

// NewEmitter create a new emitter.
// 自动注册：promauto.NewCounterVec 自动将指标注册到全局Registry，避免手动注册遗漏
func NewEmitter() *Emitter {
	return &Emitter{
		appReq: promauto.NewCounterVec(
			prometheus.CounterOpts{Subsystem: subSystem, Name: FQNAppRequest},
			TagNamesApp), // 标签顺序需匹配预定义的TagNamesApp
		realReq: promauto.NewCounterVec(
			prometheus.CounterOpts{Subsystem: subSystem, Name: FQNRealRequest},
			TagNamesReal),
		appCost: promauto.NewHistogramVec(
			prometheus.HistogramOpts{Subsystem: subSystem, Name: FQNAppCostMs},
			TagNamesApp),
		realCost: promauto.NewHistogramVec(
			prometheus.HistogramOpts{Subsystem: subSystem, Name: FQNRealCostMs},
			TagNamesReal),
	}
}

// Inc increases name by cnt with tagPairs.
func (e *Emitter) Inc(name string, cnt int, tagPairs ...string) {
	switch name {
	case FQNAppRequest:
		// tagPairs2Labels 将标签键值对转换为 prometheus.Labels，减少动态内存分配
		e.appReq.With(tagPairs2Labels(tagPairs)).Add(float64(cnt))
	case FQNRealRequest:
		e.realReq.With(tagPairs2Labels(tagPairs)).Add(float64(cnt))
	default:
	}
}

// Observe increases name by v with tagPairs.
func (e *Emitter) Observe(name string, v float64, tagPairs ...string) {
	switch name {
	case FQNAppCostMs:
		e.appCost.With(tagPairs2Labels(tagPairs)).Observe(v)
	case FQNRealCostMs:
		e.realCost.With(tagPairs2Labels(tagPairs)).Observe(v)
	default:
	}
}

// tagPairs2Labels 转换为prom的标签格式
func tagPairs2Labels(tagPairs []string) prometheus.Labels {
	labels := make(prometheus.Labels)
	for i := 0; i+1 < len(tagPairs); i += 2 {
		labels[tagPairs[i]] = tagPairs[i+1]
	}
	return labels
}
