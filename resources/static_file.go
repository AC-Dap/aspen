package resources

import (
	"aspen/router"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type StaticFile struct {
	filepath string
	router.BaseResource
}

func NewStaticFile(id string, filepath string) *StaticFile {
	return &StaticFile{
		filepath: filepath,
		BaseResource: router.BaseResource{
			Id:     id,
			Status: router.NotStarted,
		},
	}
}

func (sr *StaticFile) Start() error {
	sr.BaseResource.Status = router.Started
	return nil
}

func (sr *StaticFile) Stop() error {
	sr.BaseResource.Status = router.Stopped
	return nil
}

/*
Adds a single GET handler returning the static file under the given path.
*/
func (sr *StaticFile) AddHandlers(path string, router *httprouter.Router) error {
	router.GET(path, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		http.ServeFile(w, req, sr.filepath)
	})
	return nil
}
