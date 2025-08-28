package hystrix

import (
	metriccollector "github.com/afex/hystrix-go/hystrix/metric_collector"
)

// MetricCollectorFunc creates a MetricCollector function type definition.
type MetricCollectorFunc func(string) metriccollector.MetricCollector

// RegisterCollector is a register the statistical data collector.
// Call this function to register after implementing the collector according to its own actual business.
func RegisterCollector(collector MetricCollectorFunc) {
	metriccollector.Registry.Register(collector)
}
