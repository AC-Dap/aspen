package benchmarks

import (
	"aspen/router"
	"io"
	"log"
	"net/http"
	"testing"
)

// From https://github.com/julienschmidt/go-http-routing-benchmark
type mockResponseWriter struct{}

func (m mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m mockResponseWriter) WriteHeader(int) {}

func init() {
	log.SetOutput(io.Discard)
}

func BenchmarkRouter(b *testing.B) {
	rng := GetRNG()
	resource := &TestResource{
		BaseResource: router.BaseResource{
			Id:     "test",
			Status: router.NotStarted,
		},
	}

	// Set up router with paths
	paths := GenerateRandomPaths(rng, 1000)
	resources := make(map[string]router.Resource)
	for _, path := range paths {
		resources[path] = resource
	}
	router.UpdateRouter(resources)

	// Sample paths to get the requests we'll be benchmarking
	requests := make([]int, 10000)
	for i := range requests {
		requests[i] = rng.Intn(len(paths))
	}

	// Shared request for all requests (from https://github.com/julienschmidt/go-http-routing-benchmark)
	w := mockResponseWriter{}
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	// Benchmark how long it takes to handle every request
	b.ReportAllocs()
	for b.Loop() {
		for _, req := range requests {
			r.RequestURI = paths[req]
			u.Path = paths[req]
			u.RawQuery = rq
			router.GlobalRouter.ServeHTTP(w, r)
		}
	}
}
