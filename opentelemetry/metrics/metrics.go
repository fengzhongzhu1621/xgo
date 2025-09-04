package metrics

import "sync"

var (
	// metricsSinks emits same metrics information to multi external system at the same time.
	metricsSinksMutex = sync.RWMutex{}
	metricsSinks      = map[string]Sink{}
)

var (
	countersMutex = sync.RWMutex{}
	counters      = map[string]IMetricsCounter{}
)

// Counter creates a named counter.
func Counter(name string) IMetricsCounter {
	countersMutex.RLock()
	c, ok := counters[name]
	countersMutex.RUnlock()
	if ok && c != nil {
		return c
	}

	countersMutex.Lock()
	c, ok = counters[name]
	if ok && c != nil {
		countersMutex.Unlock()
		return c
	}
	c = &counter{name: name}
	counters[name] = c
	countersMutex.Unlock()

	return c
}
