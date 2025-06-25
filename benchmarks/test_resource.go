package benchmarks

import (
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TestResource struct {
	router.BaseResource
}

func (r *TestResource) Start() error {
	r.BaseResource.Status = router.Started
	return nil
}

func (r *TestResource) Stop() error {
	r.BaseResource.Status = router.Stopped
	return nil
}

func (r *TestResource) AddHandlers(path string, router *httprouter.Router) error {
	response := []byte("Hello World!")
	router.GET(path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		w.Write(response)
	})
	return nil
}
