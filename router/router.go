package router

import (
	"aspen/resources"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func (r *Router) Init(resources map[string]resources.Resource) {
	r.router = httprouter.New()

	for path, resource := range resources {
		resource.AddHandlers(path, r.router)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.router == nil {
		log.Fatal("Router is not initialized")
	}

	log.Printf("Request received: %s %s", req.Method, req.URL.Path)
	r.router.ServeHTTP(w, req)
}
