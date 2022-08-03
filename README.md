# promware [![Test Status](https://github.com/imbue11235/promware/workflows/Go/badge.svg)](https://github.com/imbue11235/promware/actions?query=workflow:Go) [![Go Reference](https://pkg.go.dev/badge/github.com/imbue11235/promware.svg)](https://pkg.go.dev/github.com/imbue11235/promware)

> A simple, configurable middleware for recording request metrics with [Prometheus](https://github.com/prometheus/prometheus)

## 🛠  Installation

Make sure to have Go installed (Version `1.16` or higher).

Install `promware` with `go get`:

```sh
$ go get -u github.com/imbue11235/promware
```

## 💻  Usage

### With standard library

```go
middleware := promware.Default()

handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "hello world")
})

mux := http.NewServeMux()
mux.Handle("/metrics", promhttp.Handler())
mux.Handle("/", middleware(handler))

http.ListenAndServe("<addr>", mux)
```

### Options
```go
middleware := promware.New(
    // Skip depending on request data 
    promware.WithSkipFunc(func(r *http.Request) bool {
        return r.URL.Path == "/something"
    })
    	
    // Set subsystem
    promware.WithSubsystem("my-subsystem")
    
    // Set namespace 
    promware.WithNamespace("my-namespace")
    
    // Add request counter, set name
    promware.WithRequestCounter("total-requests")
    
    // Add latency histogram, set name and buckets
    promware.WithLatencyHistogram("request-latency", []float64{0.5, 1, 2})
)
```

## 📜 License

This project is licensed under the [MIT license](LICENSE).