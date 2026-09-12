package resources

import (
	"aspen/router"
)

func RegisterResources() {
	router.RegisterResourceConstructor("static_file", NewStaticFile)
	router.RegisterResourceConstructor("directory", NewStaticDirectory)
	router.RegisterResourceConstructor("api", NewRouterAPIResource)
	router.RegisterResourceConstructor("auth", NewAuthResource)
	router.RegisterResourceConstructor("redirect", NewRedirectResource)
	router.RegisterResourceConstructor("proxy", NewProxyResource)
}
