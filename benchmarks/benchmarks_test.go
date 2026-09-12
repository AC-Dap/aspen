package benchmarks

import (
	"aspen/config"
	"aspen/logging"
	"aspen/middleware"
	"aspen/router"
	"net/http"
	"strconv"
	"testing"

	"github.com/rs/zerolog"
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
	logging.SetLoggingLevel(zerolog.Disabled)
}

func benchmarkHelper(b *testing.B, middlewareConfigs config.AllMiddlewareConfigs) {
	b.Helper()

	rng := GetRNG()
	router.RegisterResourceConstructor("test", NewTestResource)
	resourceConfig := config.ResourceConfig{
		Type:   "test",
		Params: make(map[string]any),
	}

	// Set up router with paths
	paths := GenerateRandomPaths(rng, 1000)
	routeConfigs := make([]config.RouteConfig, 1000)
	for i, path := range paths {
		routeConfigs[i] = config.RouteConfig{
			Id:          strconv.Itoa(i),
			Route:       path,
			AccessRoles: make([]string, 0),
			Resource:    resourceConfig,
		}
	}

	config := config.Config{
		Version:    config.VersionNumber{Major: config.MAJOR_VERSION, Minor: 0},
		Middleware: middlewareConfigs,
		Routes:     routeConfigs,
		Services:   make([]config.ServiceConfig, 0),
	}
	routerInstance, err := router.NewRouterInstance(config)
	if err != nil {
		b.Fatalf("error parsing config: %v", err)
	}
	router.Update(routerInstance)

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

func BenchmarkRouter(b *testing.B) {
	benchmarkHelper(b, make(config.AllMiddlewareConfigs))
}

func BenchmarkRouterWithMiddleware(b *testing.B) {
	// TODO: Ues test middleware instead
	router.RegisterMiddlewareConstructor("logger", middleware.NewLogger)

	middlewareConfigs := make(config.AllMiddlewareConfigs)
	middlewareConfigs["logger"] = config.MiddlewareConfig{
		Disabled: false,
		Params:   make(map[string]any),
		Priority: 0,
	}
	benchmarkHelper(b, middlewareConfigs)
}
