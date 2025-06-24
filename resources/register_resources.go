package resources

import "aspen/config"

func RegisterResources() {
	config.RegisterResourceConstructor("static_file", NewStaticFile)
	config.RegisterResourceConstructor("directory", NewStaticDirectory)
	config.RegisterResourceConstructor("update_router", NewUpdateRouterResource)
	config.RegisterResourceConstructor("redirect", NewRedirectResource)
}
