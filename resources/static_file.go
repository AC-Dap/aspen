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

type StaticFileParams struct {
	Filepath string
}

func NewStaticFile(base router.BaseResource, params StaticFileParams) router.Resource {
	return &StaticFile{
		filepath:     params.Filepath,
		BaseResource: base,
	}
}

/*
Adds a single GET handler returning the static file under the given path.
*/
func (sr *StaticFile) AddHandlers(path string, router *router.RouterInstance) error {
	router.GET(path, sr.BaseResource, func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		http.ServeFile(w, req, sr.filepath)
	})
	return nil
}
