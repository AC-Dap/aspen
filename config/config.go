package config

import (
	"aspen/router"
	"encoding/json"
	"fmt"
)

type Config struct {
	Middleware []MiddlewareConfig
	Routes     []RouteConfig
}

func ParseJSON(jsonBlob []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(jsonBlob, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
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
