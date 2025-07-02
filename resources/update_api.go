package resources

import (
	"aspen/router"
)

type UpdateRouterResource struct {
	router.BaseResource
}

type UpdateRouterParams struct{}

func NewUpdateRouterResource(base router.BaseResource, params UpdateRouterParams) router.Resource {
	return &UpdateRouterResource{
		BaseResource: base,
	}
}

func (ur *UpdateRouterResource) Start() error {
	ur.BaseResource.Status = router.Started
	return nil
}

func (ur *UpdateRouterResource) Stop() error {
	ur.BaseResource.Status = router.Stopped
	return nil
}

/*
Add a POST handler that listens for JSON of new resources to list.
*/
func (ur *UpdateRouterResource) AddHandlers(path string, r *router.RouterInstance) error {
	return nil
}
