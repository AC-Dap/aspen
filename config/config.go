package config

import "aspen/logging"

var lg = logging.NewTaggedLogger("ConfigSchema")

type Config struct {
	Version    VersionNumber
	Middleware AllMiddlewareConfigs
	Routes     []RouteConfig
	Services   []ServiceConfig
}
