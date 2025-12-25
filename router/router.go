package router

import (
	"aspen/logging"
	"aspen/router/service"
	"aspen/utils"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
)

var lg = logging.NewTaggedLogger("Router")

var GlobalRouter router

// Aspen router. Kept private so all instances are made through GlobalRouter.
type router struct {
	router atomic.Pointer[RouterInstance]
}

type RouterInstance struct {
	middlewareHandlers []MiddlewareHandler

	// Maps service IDs to their respective Service instances.
	services map[string]*service.Service

	// The actual HTTP router instance that handles requests.
	router *httprouter.Router
}

// Creates a new router instance with the provided middleware, services, and resources.
func NewRouterInstance(middleware []Middleware, services []*service.Service, resources map[string]Resource) *RouterInstance {
	instance := &RouterInstance{
		middlewareHandlers: make([]MiddlewareHandler, len(middleware)),
		services:           make(map[string]*service.Service),
		router:             httprouter.New(),
	}

	// Store middleware handlers
	for i, m := range middleware {
		instance.middlewareHandlers[i] = m.Handle
	}

	// Map services by their ID
	for _, service := range services {
		instance.services[service.GetID()] = service
	}

	lg.Info().Msg("Creating resource handlers for new router instance:")
	for path, resource := range resources {
		err := resource.AddHandlers(path, instance)
		if err != nil {
			lg.Warn().Str("path", path).Err(err).Msg("Error adding handlers")
		} else {
			lg.Info().Str("path", path).Str("id", resource.GetID()).Type("resource", resource).Send()
		}
	}

	return instance
}

// Initialize starts the services of the provided instance, and points
// the global router to it. Should only be called once.
func Initialize(instance *RouterInstance) error {
	lg.Info().Msg("Initializing global router instance")

	err := instance.BuildAndStartServices()
	if err != nil {
		return fmt.Errorf("error starting services: %w", err)
	}

	old := GlobalRouter.router.Swap(instance)
	if old != nil {
		return fmt.Errorf("router already has a non-nil instance")
	}

	return nil
}

// Update swaps the global router instance, and stops the old instance.
func Update(instance *RouterInstance) {
	lg.Info().Msg("Updating global router instance")
	old := GlobalRouter.router.Swap(instance)
	if old != nil {
		lg.Info().Msg("Stopping old router instance services")
		if err := old.StopServices(); err != nil {
			lg.Error().Err(err).Msg("Error stopping old router instance services")
		}
	}
}

// Shutdown stops all the currently running services and clears the router.
func Shutdown() error {
	lg.Info().Msg("Shutting down global router instance")
	router := GlobalRouter.router.Swap(nil)
	if router != nil {
		if err := router.StopServices(); err != nil {
			return fmt.Errorf("error stopping services during shutdown: %v", err)
		}
	}

	return nil
}

// ServeHTTP forwards the request to the current router instance to handle.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router := r.router.Load()
	if router == nil {
		lg.Fatal().Msg("Router is not initialized")
	}

	router.router.ServeHTTP(w, req)
}

// GetService retrieves a service by its ID from the router instance.
func (r *RouterInstance) GetService(id string) *service.Service {
	return r.services[id]
}

// BuildServices builds each service for this router instance.
func (r *RouterInstance) BuildServices() error {
	lg.Info().Msg("Building services")
	for id, service := range r.services {
		if err := service.Build(); err != nil {
			return fmt.Errorf("error building service %s: %w", id, err)
		}
	}
	return nil
}

// StartServices starts each service for this router instance.
func (r *RouterInstance) StartServices() error {
	lg.Info().Msg("Starting services")
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
	if len(r.middlewareHandlers) == 0 {
		r.router.Handle(method, path, handle)
		return
	}

	handleWithMiddleware := func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		trw := utils.NewTrackingResponseWriter(w)

		// Start recursive chain at head
		r.middlewareHandlers[0](resource, trw, req, ps, r.middlewareHandlers[1:], handle)
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
