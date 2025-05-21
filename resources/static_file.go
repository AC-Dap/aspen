package resources

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type StaticFile struct {
	id       string
	filepath string
	status   Status
}

func NewStaticFile(id string, filepath string) *StaticFile {
	return &StaticFile{
		id:       id,
		filepath: filepath,
		status:   NotStarted,
	}
}

func (sr *StaticFile) Start() error {
	sr.status = Started
	return nil
}

func (sr *StaticFile) Stop() error {
	sr.status = Stopped
	return nil
}

func (sr *StaticFile) Status() Status {
	return sr.status
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
