package metrics

type IHistogram interface {
	// Observe a histogram should have an Observe method like prometheus.
	Observe(name string, v float64, tagPairs ...string)
}

var _ IHistogram = (*histogramWrapper)(nil)

type histogramWrapper struct {
	IHistogram
	tagPairs []string
}

func WrapHistogram(h IHistogram, tagPairs ...string) *histogramWrapper {
	return &histogramWrapper{
		IHistogram: h,
		tagPairs:   tagPairs,
	}
}

// Observe implement counter.Observe.
func (h *histogramWrapper) Observe(name string, v float64, tagPairs ...string) {
	h.IHistogram.Observe(name, v, append(h.tagPairs, tagPairs...)...)
}
