package types

type Resource struct {
	Name       string
	Route      string
	Source     string
	Restricted bool
}

type ResourcesConfig struct {
	Resources []Resource
}
