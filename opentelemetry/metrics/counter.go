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
