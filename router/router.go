package router

import (
	"aspen/router/service"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

var GlobalRouter router

// Aspen router. Kept private so all instances are made through GlobalRouter.
type router struct {
	router atomic.Pointer[RouterInstance]
}

type RouterInstance struct {
	middleware []Middleware

	// Maps service IDs to their respective Service instances.
	services map[string]*service.Service

	// The actual HTTP router instance that handles requests.
	router *httprouter.Router
}

// Creates a new router instance with the provided middleware, services, and resources.
func NewRouterInstance(middleware []Middleware, services []*service.Service, resources map[string]Resource) *RouterInstance {
	instance := &RouterInstance{
		middleware: middleware,
		services:   make(map[string]*service.Service),
		router:     httprouter.New(),
	}

	// Map services by their ID
	for _, service := range services {
		instance.services[service.GetID()] = service
	}

	log.Info().Msg("Creating resource handlers for new router instance:")
	for path, resource := range resources {
		err := resource.AddHandlers(path, instance)
		if err != nil {
			log.Warn().Str("path", path).Err(err).Msg("Error adding handlers")
		} else {
			log.Info().Str("path", path).Str("id", resource.GetID()).Type("resource", resource).Send()
		}
	}

	return instance
}

// UpdateRouter swaps the global router instance, and stops the old instance.
func UpdateRouter(instance *RouterInstance) {
	log.Info().Msg("Updating global router instance")
	old := GlobalRouter.router.Swap(instance)
	if old != nil {
		log.Info().Msg("Stopping old router instance services")
		if err := old.StopServices(); err != nil {
			log.Error().Err(err).Msg("Error stopping old router instance services")
		}
	}
}

// ServeHTTP forwards the request to the current router instance to handle.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router := r.router.Load()
	if router == nil {
		log.Fatal().Msg("Router is not initialized")
	}

	router.router.ServeHTTP(w, req)
}

func (r *router) Shutdown() error {
	log.Info().Msg("Shutting down global router instance")
	router := r.router.Swap(nil)
	if router != nil {
		if err := router.StopServices(); err != nil {
			return fmt.Errorf("error stopping services during shutdown: %v", err)
		}
	}

	return nil
}

// GetService retrieves a service by its ID from the router instance.
func (r *RouterInstance) GetService(id string) *service.Service {
	return r.services[id]
}

// BuildServices builds each service for this router instance.
func (r *RouterInstance) BuildServices() error {
	for id, service := range r.services {
		if err := service.Build(); err != nil {
			return fmt.Errorf("error building service %s: %w", id, err)
		}
	}
	return nil
}

// StartServices starts each service for this router instance.
func (r *RouterInstance) StartServices() error {
	for id, service := range r.services {
		if err := service.Start(); err != nil {
			return fmt.Errorf("error starting service %s: %w", id, err)
		}
	}
	return nil
}

// BuildAndStartServices builds and starts all services managed by the router instance.
func (r *RouterInstance) BuildAndStartServices() error {
	err := r.BuildServices()
	if err == nil {
		err = r.StartServices()
	}
	return err
}

// StopServices calls Stop() on all services managed by the router instance.
func (r *RouterInstance) StopServices() error {
	for id, service := range r.services {
		if err := service.Stop(); err != nil {
			return fmt.Errorf("error stopping service %s: %w", id, err)
		}
	}
	return nil
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
