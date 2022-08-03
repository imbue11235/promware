package promware

import (
	"fmt"
	"github.com/appleboy/gofight/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	middleware := Default()

	routes := []struct {
		path    string
		handler http.Handler
	}{
		{
			path:    "/metrics",
			handler: promhttp.Handler(),
		},
		{
			path: "/",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "hello world")
			}),
		},
		{
			path: "/other",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "no")
			}),
		},
	}

	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.path, middleware(route.handler))
	}

	r := gofight.New()
	r.GET("/").Run(mux, func(rsp gofight.HTTPResponse, req gofight.HTTPRequest) {
		assert.Equal(t, "hello world", rsp.Body.String())
		assert.Equal(t, http.StatusOK, rsp.Code)
	})

	r.GET("/metrics").Run(mux, func(rsp gofight.HTTPResponse, req gofight.HTTPRequest) {
		assert.Equal(t, http.StatusOK, rsp.Code)
		body := rsp.Body.String()

		assert.Contains(t, body, "request_duration_seconds")
		assert.Contains(t, body, "requests_total{code=\"200\",method=\"GET\",url=\"/\"} 1")
		// metrics endpoints should be skipped
		assert.NotContains(t, body, "requests_total{code=\"200\",method=\"GET\",url=\"/metrics\"} 1")
		// other was never called, and should be skipped
		assert.NotContains(t, body, "requests_total{code=\"200\",method=\"GET\",url=\"/other\"} 1")
	})

}
