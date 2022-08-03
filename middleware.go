package promware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Middleware is responsible for creating a middleware handler func
type Middleware func(next http.Handler) http.Handler

// New creates a new prometheus middleware from given
// options.
// Note: the middleware has no prometheus collectors attached
// per default. If you want the default collectors, use Default
func New(options ...option) Middleware {
	col := newCollector()
	col.apply(options...)

	if err := col.register(); err != nil {
		log.Println("could not register collectors", err)

		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// create a delegate response writer
			rw := newResponseWriter(w)

			// call next handler
			next.ServeHTTP(rw, r)

			if col.shouldSkip(r) {
				return
			}

			// record the elapsed time
			elapsed := float64(time.Since(start)) / float64(time.Second)

			// get the URL path
			path := r.URL.Path

			// ensure method is uppercase, e.g. "GET", "POST" etc.
			method := strings.ToUpper(r.Method)

			// convert HTTP status code to a string for usage with prometheus
			code := strconv.Itoa(rw.Status())

			go col.incrementRequestCounter(code, method, path)
			go col.addObservationToLatencyHistogram(code, method, path, elapsed)
		})
	}
}

// Default creates a new middleware with default collectors attached
func Default() Middleware {
	return New(
		WithDefaultLatencyHistogram(),
		WithDefaultRequestCounter(),
	)
}
