package benchmarks

import (
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TestResource struct {
	router.BaseResource
}

type TestResourceParams struct{}

func NewTestResource(base router.BaseResource, params TestResourceParams) router.Resource {
	return &TestResource{
		BaseResource: base,
	}
}

func (r *TestResource) AddHandlers(path string, router *router.RouterInstance) error {
	response := []byte("Hello World!")
	router.GET(path, r.BaseResource, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Write(response)
	})
	return nil
}
