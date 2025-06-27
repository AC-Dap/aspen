package resources

import (
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
func (ur *UpdateRouterResource) AddHandlers(path string, r *httprouter.Router) error {
	r.POST(path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var resources = map[string]router.Resource{
			"/info":   NewStaticFile("info", "README.md"),
			"/design": NewStaticFile("design", "design.md"),
			"/code":   NewStaticDirectory("resources", "resources", []string{"static_file.go"}, false),
			"/go.mod": NewStaticFile("go_mod", "go.mod"),
		}

		router.UpdateRouter(resources)
	})
	return nil
}
