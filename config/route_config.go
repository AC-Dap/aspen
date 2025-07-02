package config

import (
	"aspen/router"
	"fmt"
)

type RouteConfig struct {
	Route    string
	Id       string
	Resource ResourceConfig
}

func (rc RouteConfig) Parse() (router.Resource, error) {
	// Create base resource
	base := router.BaseResource{Id: rc.Id, Status: router.NotStarted}

	// Parse resource
	newResource, err := rc.Resource.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("error parsing \"%s\" route: %w", rc.Id, err)
	}

	return newResource, nil
}
