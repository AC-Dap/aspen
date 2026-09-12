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
	Version string `json:"timestamp"`
}

// verifyVersion checks if the given version matches the given config version.
func verifyVersion(r *http.Request, c *config.Config) error {
	var params postParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return fmt.Errorf("error parsing post params: %w", err)
	}

	var versionNumber config.VersionNumber
	if err := versionNumber.UnmarshalText([]byte(params.Version)); err != nil {
		return fmt.Errorf("error parsing version from post params: %w", err)
	}

	if versionNumber != c.Version {
		return fmt.Errorf("mismatching version numbers: expected %v, got %v", c.Version, versionNumber)
	}

	return nil
}

// verifyMajorVersion checks if the given version has the correct major version.
func verifyMajorVersion(r *http.Request) error {
	var params postParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return fmt.Errorf("error parsing post params: %w", err)
	}

	var versionNumber config.VersionNumber
	if err := versionNumber.UnmarshalText([]byte(params.Version)); err != nil {
		return fmt.Errorf("error parsing version from post params: %w", err)
	}

	if !versionNumber.IsCompatibleVersion() {
		return fmt.Errorf("mismatching major version numbers: expected %v, got %v", config.MAJOR_VERSION, versionNumber.Major)
	}

	return nil
}

// AddHandlers adds API routes that allow querying and updating the router config.
func (ur *RouterAPIResource) AddHandlers(path string, r *router.RouterInstance) error {
	/*
		 		* GET version: router config version string
			 	* GET middleware: Array of strings
				* GET routes: Array of route JSONs
				* GET route(id): Route JSON corresponding to given id
				* GET services: Array of service JSONs
				* GET service(id): Service JSON corresponding to given id

				* GET available_middleware: Array of all middleware strings
				* GET available_resources: Array of resource type strings
				* GET resource_params(type): Return params for the given resource type

				- Each POST request should also include the config version number to prevent replay attacks
				* POST set_middleware(middleware): Sets middleware to be new array of strings
				* POST add_route(route): Adds a new route
				* POST delete_route(id): Deletes the route with the given id
				* POST update_route(id, resource): Updates the route resource with the given id
				* POST change_route(id, route): Changes the route path for the given id

				TODO: remove this once we have a file watcher autoreloading
				* POST reload: Reloads the router config from disk
	*/
	r.GET(path+"/version", ur.BaseResource, gen_get_version(r))
	r.GET(path+"/middleware", ur.BaseResource, gen_get_middleware(r))
	r.GET(path+"/routes", ur.BaseResource, gen_get_routes(r))
	r.GET(path+"/route/:id", ur.BaseResource, gen_get_route(r))

	r.GET(path+"/available_middleware", ur.BaseResource, gen_get_available_middleware(r))
	r.GET(path+"/available_resources", ur.BaseResource, gen_get_available_resources(r))
	r.GET(path+"/resource_params/:type", ur.BaseResource, gen_get_resource_params(r))

	r.POST(path+"/set_middleware", ur.BaseResource, gen_set_middleware(r))
	r.POST(path+"/add_route", ur.BaseResource, gen_add_route(r))
	r.POST(path+"/delete_route", ur.BaseResource, gen_delete_route(r))
	r.POST(path+"/update_route", ur.BaseResource, gen_update_route(r))
	r.POST(path+"/change_route", ur.BaseResource, gen_change_route(r))

	r.POST(path+"/reload", ur.BaseResource, reload)

	return nil
}

/*
 * ===========================================================
 * Below are the handlers for the API endpoints
 * ===========================================================
 */

func gen_get_version(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		versionString, err := ri.Config.Version.MarshalText()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal config version: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(versionString)
	}
}

func gen_get_middleware(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(ri.Config.Middleware)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func gen_get_routes(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(ri.Config.Routes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func gen_get_route(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		routeID := p.ByName("id")
		for _, route := range ri.Config.Routes {
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
}

func gen_get_available_middleware(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyMajorVersion(r)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(config.AvailableMiddleware())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func gen_get_available_resources(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyMajorVersion(r)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(config.AvailableResources())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func gen_get_resource_params(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyMajorVersion(r)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

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
}

func gen_set_middleware(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		var body struct {
			Middleware config.AllMiddlewareConfigs `json:"middleware"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		err = router.UpdateGlobalConfig(func(config *config.Config) error {
			config.Middleware = body.Middleware
			return nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update middleware: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func gen_add_route(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		var body struct {
			Route config.RouteConfig `json:"route"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		err = router.UpdateGlobalConfig(func(config *config.Config) error {
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
}

func gen_delete_route(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		var body struct {
			Id string `json:"id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		err = router.UpdateGlobalConfig(func(config *config.Config) error {
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
}

func gen_update_route(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		var body struct {
			Id       string                `json:"id"`
			Resource config.ResourceConfig `json:"resource"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		err = router.UpdateGlobalConfig(func(config *config.Config) error {
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
}

func gen_change_route(ri *router.RouterInstance) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := verifyVersion(r, &ri.Config)
		if err == nil {
			http.Error(w, fmt.Sprintf("Invalid config version: %v", err), http.StatusBadRequest)
			return
		}

		var body struct {
			Id    string `json:"id"`
			Route string `json:"route"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode body: %v", err), http.StatusBadRequest)
			return
		}

		err = router.UpdateGlobalConfig(func(config *config.Config) error {
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
}

func reload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Load config
	instance, err := router.ParseGlobalConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse global config: %v", err), http.StatusInternalServerError)
		return
	}

	// Init router
	router.Update(instance)

	w.WriteHeader(http.StatusOK)
}
