package resources

import (
	"aspen/router"
)

type UpdateRouterResource struct {
	router.BaseResource
}

func NewUpdateRouterResource(id string) *UpdateRouterResource {
	return &UpdateRouterResource{
		BaseResource: router.BaseResource{
			Id:     id,
			Status: router.NotStarted,
		},
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
