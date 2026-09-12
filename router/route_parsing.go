package router

import (
	"aspen/config"
	"fmt"
)

func ParseRoute(rc config.RouteConfig) (Resource, error) {
	// Create base resource
	base := NewBaseResource(rc.Id, rc.AccessRoles)

	// Parse resource
	newResource, err := ParseResource(base, rc.Resource)
	if err != nil {
		return nil, fmt.Errorf("error parsing \"%s\" route: %w", rc.Id, err)
	}

	return newResource, nil
}
