package router

import (
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
func UpdateRouter(resources map[string]Resource) {
	router := httprouter.New()

	log.Println("Updating router with resources:")
	for path, resource := range resources {
		err := resource.AddHandlers(path, router)
		if err != nil {
			log.Printf("  %s ⚠ Error adding handlers: %v", path, err)
		} else {
			log.Printf("  %s -> %s (%T)", path, resource.GetID(), resource)
		}
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
