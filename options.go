package promware

import "github.com/prometheus/client_golang/prometheus"

// option alters the recorder struct
type option func(c *collector)

func WithSkipFunc(skipFunc SkipFunc) option {
	return func(c *collector) {
		c.shouldSkip = skipFunc
	}
}

// WithNamespace sets the namespace of the prometheus collectors
func WithNamespace(namespace string) option {
	return func(c *collector) {
		c.namespace = namespace
	}
}

// WithSubsystem sets the subsystem of the prometheus collectors
func WithSubsystem(name string) option {
	return func(c *collector) {
		c.subsystem = name
	}
}

// WithRequestCounter creates a new request counter collector with a given name
func WithRequestCounter(name string) option {
	return func(c *collector) {
		c.requestCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: c.namespace,
				Subsystem: c.subsystem,
				Name:      name,
				Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
			},
			[]string{"code", "method", "url"},
		)
	}
}

// WithDefaultRequestCounter creates a request counter with the default name `requests_total`
func WithDefaultRequestCounter() option {
	return WithRequestCounter("requests_total")
}

// WithLatencyHistogram creates a latency histogram collector with a given name
func WithLatencyHistogram(name string, buckets []float64) option {
	return func(c *collector) {
		c.latencyHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: c.namespace,
				Subsystem: c.subsystem,
				Name:      name,
				Help:      "How long it took to process the request in seconds, partitioned by status code, method and HTTP path.",
				Buckets:   buckets,
			},
			[]string{"code", "method", "url"},
		)
	}
}

// WithDefaultLatencyHistogram creates a latency histogram collector with the name `request_duration_seconds`
func WithDefaultLatencyHistogram() option {
	return WithLatencyHistogram("request_duration_seconds", prometheus.DefBuckets)
}
