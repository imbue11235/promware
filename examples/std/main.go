package main

import (
	"fmt"
	"github.com/imbue11235/promware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	addr := ":8080"

	middleware := promware.Default()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "hello world")
	})

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", middleware(handler))

	log.Printf("Serving on addr %s", addr)

	http.ListenAndServe(addr, mux)
}
