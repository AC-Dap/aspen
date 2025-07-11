package resources

import (
	"aspen/config"
	"aspen/router"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RouterAPIResource struct {
	router.BaseResource
}

type RouterAPIParams struct{}

func NewRouterAPIResource(base router.BaseResource, params RouterAPIParams) router.Resource {
	return &RouterAPIResource{
		BaseResource: base,
	}
}

// postParams is expected in POST requests to verify the timestamp
type postParams struct {
	Timestamp int64 `json:"timestamp"`
}

// verifyTimestamp checks if the given timestamp is before the last updated time of the config.
// If things are OK, it updates the config's LastUpdated field to the given timestamp.
func verifyTimestamp(p postParams, config *config.Config) error {
	if config.LastUpdated > p.Timestamp {
		return fmt.Errorf("given timestamp is in the past")
	}
	config.LastUpdated = p.Timestamp
	return nil
}

// Adds API routes that allow querying and updating the router config.
func (ur *RouterAPIResource) AddHandlers(path string, r *router.RouterInstance) error {
	/*
		 	* GET middleware: Array of strings
			* GET routes: Array of route JSONs
			* GET route(id): Route JSON corresponding to given id
			* GET services: Array of service JSONs
			* GET service(id): Service JSON corresponding to given id

			* GET available_middleware: Array of all middleware strings
			* GET available_resources: Array of resource type strings
			* GET resource_params(type): Return params for the given resource type

			- Each POST request should also include a timestamp field to prevent replay attacks
			* POST set_middleware(middleware): Sets middleware to be new array of strings
			* POST add_route(route): Adds a new route
			* POST delete_route(id): Deletes the route with the given id
			* POST update_route(id, resource): Updates the route resource with the given id
			* POST change_route(id, route): Changes the route path for the given id

			* POST reload: Reloads the router config from disk
	*/
	r.GET(path+"/middleware", ur.BaseResource, get_middleware)
	r.GET(path+"/routes", ur.BaseResource, get_routes)
	r.GET(path+"/route/:id", ur.BaseResource, get_route)

	r.GET(path+"/available_middleware", ur.BaseResource, get_available_middleware)
	r.GET(path+"/available_resources", ur.BaseResource, get_available_resources)
	r.GET(path+"/resource_params/:type", ur.BaseResource, get_resource_params)

	r.POST(path+"/set_middleware", ur.BaseResource, set_middleware)
	r.POST(path+"/add_route", ur.BaseResource, add_route)
	r.POST(path+"/delete_route", ur.BaseResource, delete_route)
	r.POST(path+"/update_route", ur.BaseResource, update_route)
	r.POST(path+"/change_route", ur.BaseResource, change_route)

	r.POST(path+"/reload", ur.BaseResource, reload)

	return nil
}

/*
 * ===========================================================
 * Below are the handlers for the API endpoints
 * ===========================================================
 */

func get_middleware(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := config.ReadGlobalConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read global config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(c.Middleware)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func get_routes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := config.ReadGlobalConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read global config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(c.Routes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func get_route(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := config.ReadGlobalConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read global config: %v", err), http.StatusInternalServerError)
		return
	}

	routeID := p.ByName("id")
	for _, route := range c.Routes {
		if route.Id == routeID {
			w.Header().Set("Content-Type", "application/json")
			data, err := json.Marshal(route)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		}
	}
	http.Error(w, "Route not found", http.StatusNotFound)
}

func get_available_middleware(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	middleware := config.AvailableMiddleware()

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(middleware)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func get_available_resources(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resources := config.AvailableResources()

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resources)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func get_resource_params(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resType := p.ByName("type")
	resource_params, err := config.GetResourceParams(resType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get resource params: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resource_params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func set_middleware(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Middleware []config.MiddlewareConfig `json:"middleware"`
		postParams
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	err := config.UpdateGlobalConfig(func(config *config.Config) error {
		if err := verifyTimestamp(body.postParams, config); err != nil {
			return err
		}

		config.Middleware = body.Middleware
		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update middleware: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func add_route(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Route config.RouteConfig `json:"route"`
		postParams
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	err := config.UpdateGlobalConfig(func(config *config.Config) error {
		if err := verifyTimestamp(body.postParams, config); err != nil {
			return err
		}

		// Make sure the route ID is unique
		for _, route := range config.Routes {
			if route.Id == body.Route.Id {
				return fmt.Errorf("route with ID %s already exists", body.Route.Id)
			}
		}

		config.Routes = append(config.Routes, body.Route)
		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add route: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func delete_route(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Id string `json:"id"`
		postParams
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	err := config.UpdateGlobalConfig(func(config *config.Config) error {
		if err := verifyTimestamp(body.postParams, config); err != nil {
			return err
		}

		for i, route := range config.Routes {
			if route.Id == body.Id {
				// Remove the route by replacing it with the last element and slicing
				config.Routes[i] = config.Routes[len(config.Routes)-1]
				config.Routes = config.Routes[:len(config.Routes)-1]
			}
		}
		// Deleting a route that doesn't exist is not an error, just a no-op
		return nil
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete route: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func update_route(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Id       string                `json:"id"`
		Resource config.ResourceConfig `json:"resource"`
		postParams
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	err := config.UpdateGlobalConfig(func(config *config.Config) error {
		if err := verifyTimestamp(body.postParams, config); err != nil {
			return err
		}

		for i, route := range config.Routes {
			if route.Id == body.Id {
				config.Routes[i].Resource = body.Resource
				return nil
			}
		}
		return fmt.Errorf("route with ID %s doesn't exist", body.Id)
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update route: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func change_route(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var body struct {
		Id    string `json:"id"`
		Route string `json:"route"`
		postParams
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	err := config.UpdateGlobalConfig(func(config *config.Config) error {
		if err := verifyTimestamp(body.postParams, config); err != nil {
			return err
		}

		for i, route := range config.Routes {
			if route.Id == body.Id {
				config.Routes[i].Route = body.Route
				return nil
			}
		}
		return fmt.Errorf("route with ID %s doesn't exist", body.Id)
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to change route: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func reload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Load config
	instance, err := config.ParseGlobalConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse global config: %v", err), http.StatusInternalServerError)
		return
	}

	// Init router
	router.UpdateRouter(instance)

	w.WriteHeader(http.StatusOK)
}
