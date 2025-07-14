package config

import (
	"aspen/router"
	"aspen/router/service"
	"fmt"
)

type Config struct {
	LastUpdated int64
	Middleware  []MiddlewareConfig
	Routes      []RouteConfig
}

func (c *Config) GetMiddleware() ([]router.Middleware, error) {
	var middlewares = make([]router.Middleware, len(c.Middleware))
	for i, middleware := range c.Middleware {
		mw, err := middleware.Parse()
		if err != nil {
			return nil, fmt.Errorf("unable to parse middleware: %w", err)
		}
		middlewares[i] = mw
	}

	return middlewares, nil
}

func (c *Config) GetResourceRoutes() (map[string]router.Resource, error) {
	var resource_routes = make(map[string]router.Resource)
	for _, route := range c.Routes {
		resource, err := route.Parse()
		if err != nil {
			return nil, fmt.Errorf("unable to parse route: %w", err)
		}
		resource_routes[route.Route] = resource
	}

	return resource_routes, nil
}

func (c *Config) ToRouterInstance() (*router.RouterInstance, error) {
	middleware, err := c.GetMiddleware()
	if err != nil {
		return nil, fmt.Errorf("error loading middleware: %w", err)
	}

	resource_routes, err := c.GetResourceRoutes()
	if err != nil {
		return nil, fmt.Errorf("error loading routes: %w", err)
	}

	return router.NewRouterInstance(
		middleware,
		[]*service.Service{},
		resource_routes,
	), nil
}
