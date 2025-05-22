package router

import (
	"aspen/resources"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
)

var GlobalRouter router

// Aspen router. Kept private so all instances are made through GlobalRouter.
type router struct {
	router atomic.Pointer[httprouter.Router]
}

/*
Creates a new router for the given set of resources and atomically swaps the current router for this new one.
The map should map paths to resources.
*/
func UpdateRouter(resources map[string]resources.Resource) {
	router := httprouter.New()

	for path, resource := range resources {
		resource.AddHandlers(path, router)
	}
	GlobalRouter.router.Swap(router)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router := r.router.Load()
	if router == nil {
		log.Fatal("Router is not initialized")
	}

	log.Printf("Request received: %s %s", req.Method, req.URL.Path)
	router.ServeHTTP(w, req)
}
