package promware

import "net/http"

// responseWriter is an extension of the http.ResponseWriter which
// enables us to save the statusCode written from a handler func further
// down the chain
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// WriteHeader saves the status code for later usage
// and calls the original ResponseWriter's WriteHeader method
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Status returns the status code
func (w *responseWriter) Status() int {
	return w.statusCode
}
