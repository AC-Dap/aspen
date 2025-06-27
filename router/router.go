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
	router atomic.Pointer[RouterInstance]
}

type RouterInstance struct {
	middleware []Middleware
	router     *httprouter.Router
}

func NewRouterInstance(middleware []Middleware, resources map[string]Resource) *RouterInstance {
	instance := &RouterInstance{
		middleware: middleware,
		router:     httprouter.New(),
	}

	log.Println("Initializing router instance with resources:")
	for path, resource := range resources {
		err := resource.AddHandlers(path, instance)
		if err != nil {
			log.Printf("  %s ⚠ Error adding handlers: %v", path, err)
		} else {
			log.Printf("  %s -> %s (%T)", path, resource.GetID(), resource)
		}
	}

	return instance
}

// UpdateRouter swaps the global router instance.
func UpdateRouter(instance *RouterInstance) {
	log.Println("Updating global router instance")
	GlobalRouter.router.Swap(instance)
}

// ServeHTTP forwards the request to the current router instance to handle.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router := r.router.Load()
	if router == nil {
		log.Fatal("Router is not initialized")
	}

	log.Printf("Request received: %s %s", req.Method, req.URL.Path)
	router.router.ServeHTTP(w, req)
}

// Handle assigns a resource and handler to a specific method and path.
func (r *RouterInstance) Handle(method, path string, resource BaseResource, handle httprouter.Handle) {
	handleWithMiddleware := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Execute middleware in the order they were added
		for _, middleware := range r.middleware {
			if err, err_code := middleware.Handle(resource, w, req, ps); err != nil {
				http.Error(w, err.Error(), err_code)
				return
			}
		}

		// Call the resource handler
		handle(w, req, ps)
	}

	r.router.Handle(method, path, handleWithMiddleware)
}

// GET wraps the Handle method for GET requests.
func (r *RouterInstance) GET(path string, resource BaseResource, handle httprouter.Handle) {
	r.Handle(http.MethodGet, path, resource, handle)
}

// POST wraps the Handle method for POST requests.
func (r *RouterInstance) POST(path string, resource BaseResource, handle httprouter.Handle) {
	r.Handle(http.MethodPost, path, resource, handle)
}
