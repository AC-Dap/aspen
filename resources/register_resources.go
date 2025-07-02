package resources

import (
	"aspen/config"
)

func RegisterResources() {
	config.RegisterResourceConstructor[StaticFileParams]("static_file", NewStaticFile)
	config.RegisterResourceConstructor[StaticDirectoryParams]("directory", NewStaticDirectory)
	config.RegisterResourceConstructor[UpdateRouterParams]("update_router", NewUpdateRouterResource)
	config.RegisterResourceConstructor[RedirectParams]("redirect", NewRedirectResource)
	config.RegisterResourceConstructor[ProxyParams]("proxy", NewProxyResource)
}
