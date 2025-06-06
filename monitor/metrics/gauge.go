package metrics

import (
	"math"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

// Gauge TODO
type Gauge struct {
	valGauge prometheus.GaugeFunc
	maxGauge prometheus.GaugeFunc

	descC    chan<- *prometheus.Desc
	collectC chan<- *prometheus.Metric

	val uint64
	max uint64
}

// NewGauge TODO
func NewGauge(opt prometheus.GaugeOpts) *Gauge {
	g := Gauge{}
	g.valGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: opt.Name,
			Help: opt.Help,
		},
		func() float64 { return math.Float64frombits(atomic.LoadUint64(&g.val)) },
	)

	g.maxGauge = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: opt.Name + "_max",
			Help: "max number of " + opt.Name,
		},
		func() float64 { return math.Float64frombits(atomic.LoadUint64(&g.max)) },
	)

	return &g
}

// Describe TODO
func (g *Gauge) Describe(ch chan<- *prometheus.Desc) {
	g.valGauge.Describe(ch)
	g.maxGauge.Describe(ch)
}

// Collect TODO
func (g *Gauge) Collect(ch chan<- prometheus.Metric) {
	g.valGauge.Collect(ch)
	g.maxGauge.Collect(ch)
}

// Inc TODO
func (g *Gauge) Inc() {
	new := g.Add(1)
	old := atomic.LoadUint64(&g.max)
	if new > old {
		atomic.CompareAndSwapUint64(&g.max, old, new)
	}
}

// Dec TODO
func (g *Gauge) Dec() {
	g.Add(-1)
}

// Add TODO
func (g *Gauge) Add(val float64) uint64 {
	for {
		oldBits := atomic.LoadUint64(&g.val)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&g.val, oldBits, newBits) {
			return newBits
		}
	}
}
