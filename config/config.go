package config

import (
	"aspen/router"
	"encoding/json"
	"fmt"
)

type Config struct {
	Middleware []string
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
