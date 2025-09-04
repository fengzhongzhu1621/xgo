package metrics

type ICounter interface {
	// Inc a counter should have a Inc method.
	Inc(name string, cnt int, tagPairs ...string)
}

var _ ICounter = (*counterWrapper)(nil)

type counterWrapper struct {
	ICounter
	tagPairs []string
}

func WrapCounter(c ICounter, tagPairs ...string) *counterWrapper {
	return &counterWrapper{
		ICounter: c,
		tagPairs: tagPairs,
	}
}

// Inc implement counter.Inc.
func (c *counterWrapper) Inc(name string, cnt int, tagPairs ...string) {
	c.ICounter.Inc(name, cnt, append(c.tagPairs, tagPairs...)...)
}

// ICounter is the interface that emits counter type metrics.
type IMetricsCounter interface {
	// Incr increments the counter by one.
	Incr()

	// IncrBy increments the counter by delta.
	IncrBy(delta float64)
}

// counter defines the counter. counter is report to each external Sink-able system.
type counter struct {
	name string
}

// Incr increases counter by one.
func (c *counter) Incr() {
	c.IncrBy(1)
}

// IncrBy increases counter by v and reports for each external Sink-able systems.
func (c *counter) IncrBy(v float64) {
	if len(metricsSinks) == 0 {
		return
	}
	rec := NewSingleDimensionMetrics(c.name, v, PolicySUM)
	for _, sink := range metricsSinks {
		sink.Report(rec)
	}
}
