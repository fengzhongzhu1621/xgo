package rest

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ClientInterface TODO
type ClientInterface interface {
	Verb(verb VerbType) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
	Patch() *Request
}

// NewRESTClient TODO
func NewRESTClient(c *Capability, baseUrl string) ClientInterface {
	if baseUrl != "/" {
		baseUrl = strings.Trim(baseUrl, "/")
		baseUrl = "/" + baseUrl + "/"
	}

	if c.ToleranceLatencyTime <= 0 {
		// set default tolerance latency time
		c.ToleranceLatencyTime = 2 * time.Second
	}

	client := &RESTClient{
		baseUrl:    baseUrl,
		capability: c,
	}

	if c.MetricOpts.Register != nil {

		var buckets []float64
		if len(c.MetricOpts.DurationBuckets) == 0 {
			// set default buckets
			buckets = []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000}
		} else {
			// use user defined buckets
			buckets = c.MetricOpts.DurationBuckets
		}

		client.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "cmdb_apimachinary_requests_duration_millisecond",
			Help:    "third party api request duration millisecond.",
			Buckets: buckets,
		}, []string{"handler", "status_code", "dimension"})

		if err := c.MetricOpts.Register.Register(client.requestDuration); err != nil {
			if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
				client.requestDuration = are.ExistingCollector.(*prometheus.HistogramVec)
			} else {
				panic(err)
			}
		}
	}

	return client
}

// RESTClient TODO
type RESTClient struct {
	baseUrl    string
	capability *Capability

	requestDuration *prometheus.HistogramVec
}

// Verb TODO
func (r *RESTClient) Verb(verb VerbType) *Request {
	return &Request{
		parent:     r,
		verb:       verb,
		baseURL:    r.baseUrl,
		capability: r.capability,
	}
}

// Post TODO
func (r *RESTClient) Post() *Request {
	return r.Verb(POST)
}

// Put TODO
func (r *RESTClient) Put() *Request {
	return r.Verb(PUT)
}

// Get TODO
func (r *RESTClient) Get() *Request {
	return r.Verb(GET)
}

// Delete TODO
func (r *RESTClient) Delete() *Request {
	return r.Verb(DELETE)
}

// Patch TODO
func (r *RESTClient) Patch() *Request {
	return r.Verb(PATCH)
}
