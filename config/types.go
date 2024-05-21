package config

type Resource struct {
	Name       string
	Route      string
	Source     string
	Restricted bool
}

type ConfigTOML struct {
	Resources []Resource
}
