package promware

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

const defaultMetricsPath = "/metrics"

// SkipFunc allows for certain requests to be skipped
// based on information obtained from the http.Request, e.g. URLPath.
type SkipFunc func(r *http.Request) bool

// collector holds information and collectors needed by prometheus
type collector struct {
	namespace string
	subsystem string

	requestCounter   *prometheus.CounterVec
	latencyHistogram *prometheus.HistogramVec

	shouldSkip SkipFunc
}

// newCollector creates a new collector with empty values
func newCollector() *collector {
	return &collector{
		namespace: "",
		subsystem: "",
		shouldSkip: func(r *http.Request) bool {
			// default skips on `/metrics` endpoint, which
			// is commonly used with Prometheus
			return r.URL.Path == defaultMetricsPath
		},
	}
}

// apply adds options to the recorder
func (c *collector) apply(options ...option) {
	for _, opt := range options {
		opt(c)
	}
}

// register registers all the prometheus.Collectors
func (c *collector) register() error {
	collectors := []prometheus.Collector{
		c.requestCounter,
		c.latencyHistogram,
	}

	for _, collector := range collectors {
		if collector == nil {
			continue
		}

		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

// addObservationToLatencyHistogram adds an observation to the latency histogram collector
func (c *collector) addObservationToLatencyHistogram(code, method, path string, seconds float64) {
	if c.latencyHistogram == nil {
		return
	}

	c.latencyHistogram.WithLabelValues(code, method, path).Observe(seconds)
}

// incrementRequestCounter increments the current request counter collector by 1
func (c *collector) incrementRequestCounter(code, method, path string) {
	if c.requestCounter == nil {
		return
	}

	c.requestCounter.WithLabelValues(code, method, path).Inc()
}
