package metrics

// Policy is the metrics aggregation policy.
type Policy int

// All available Policy(s).
const (
	PolicyNONE      = 0 // Undefined
	PolicySET       = 1 // instantaneous value
	PolicySUM       = 2 // summary
	PolicyAVG       = 3 // average
	PolicyMAX       = 4 // maximum
	PolicyMIN       = 5 // minimum
	PolicyMID       = 6 // median
	PolicyTimer     = 7 // timer
	PolicyHistogram = 8 // histogram
)

// Sink defines the interface an external monitor system should provide.
type Sink interface {
	// Name returns the name of the monitor system.
	Name() string
	// Report reports a record to monitor system.
	Report(rec Record, opts ...Option) error
}

// Record is the single record.
//
// terminologies:
//   - dimension name
//     is an attribute of a data, often used to filter data, such as a photo album business module
//     includes region and server room.
//   - dimension value
//     refines the dimension. For example, the regions of the album business module include Shenzhen,
//     Shanghai, etc., the region is the dimension, and Shenzhen and Shanghai are the dimension
//     values.
//   - metric
//     is a measurement, used to aggregate and calculate. For example, request count of album business
//     module in ShenZhen Telecom is a metric.
type Record struct {
	Name string // the name of the record
	// dimension name: such as region, host and disk number.
	// dimension value: such as region=ShangHai.
	dimensions []*Dimension
	metrics    []*Metrics
}

// Dimension defines the dimension.
type Dimension struct {
	Name  string
	Value string
}

// Metrics defines the metric.
type Metrics struct {
	name   string  // metric name
	value  float64 // metric value
	policy Policy  // aggregation policy
}

// NewSingleDimensionMetrics creates a Record with no dimension and only one metric.
func NewSingleDimensionMetrics(name string, value float64, policy Policy) Record {
	return Record{
		dimensions: nil,
		metrics: []*Metrics{
			{name: name, value: value, policy: policy},
		},
	}
}
