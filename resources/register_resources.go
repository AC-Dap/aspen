package resources

import (
	"aspen/config"
)

func RegisterResources() {
	config.RegisterResourceConstructor[StaticFileParams]("static_file", NewStaticFile)
	config.RegisterResourceConstructor[StaticDirectoryParams]("directory", NewStaticDirectory)
	config.RegisterResourceConstructor[RouterAPIParams]("api", NewRouterAPIResource)
	config.RegisterResourceConstructor[RedirectParams]("redirect", NewRedirectResource)
	config.RegisterResourceConstructor[ProxyParams]("proxy", NewProxyResource)
}
