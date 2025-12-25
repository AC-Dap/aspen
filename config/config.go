package config

import (
	"aspen/logging"
	"aspen/router"
	"aspen/router/service"
	"fmt"
)

var lg = logging.NewTaggedLogger("Config")

type Config struct {
	LastUpdated int64
	Middleware  AllMiddlewareConfigs
	Routes      []RouteConfig
	Services    []ServiceConfig
}

func (c *Config) GetMiddleware() ([]router.Middleware, error) {
	return c.Middleware.Parse()
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

func (c *Config) GetServices() ([]*service.Service, error) {
	var services = make([]*service.Service, len(c.Services))
	for i, serviceConfig := range c.Services {
		svc, err := serviceConfig.Parse()
		if err != nil {
			return nil, fmt.Errorf("unable to parse service: %w", err)
		}
		services[i] = svc
	}

	return services, nil
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

	services, err := c.GetServices()
	if err != nil {
		return nil, fmt.Errorf("error loading services: %w", err)
	}

	return router.NewRouterInstance(
		middleware,
		services,
		resource_routes,
	), nil
}
