package resources

import (
	"aspen/config"
)

func RegisterResources() {
	config.RegisterResourceConstructor("static_file", NewStaticFile)
	config.RegisterResourceConstructor("directory", NewStaticDirectory)
	config.RegisterResourceConstructor("api", NewRouterAPIResource)
	config.RegisterResourceConstructor("auth", NewAuthResource)
	config.RegisterResourceConstructor("redirect", NewRedirectResource)
	config.RegisterResourceConstructor("proxy", NewProxyResource)
}
