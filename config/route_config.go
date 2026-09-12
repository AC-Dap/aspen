package config

type RouteConfig struct {
	Id          string
	Route       string
	AccessRoles []string
	Resource    ResourceConfig
}
