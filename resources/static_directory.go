package resources

import (
	"aspen/router"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type StaticDirectory struct {
	path string

	// Allowed files to serve within this directory.
	// Paths are relative to the base path, and '*' can be used to serve all files within a directory.
	whitelist                []string
	allow_directory_browsing bool
	router.BaseResource
}

func NewStaticDirectory(id string, path string, whitelist []string, allow_directory_browsing bool) *StaticDirectory {
	// Include everything if we have an empty whitelist
	if len(whitelist) == 0 {
		whitelist = []string{"*"}
	}

	return &StaticDirectory{
		path:                     path,
		whitelist:                whitelist,
		allow_directory_browsing: allow_directory_browsing,
		BaseResource: router.BaseResource{
			Id:     id,
			Status: router.NotStarted,
		},
	}
}

func (sd *StaticDirectory) Start() error {
	sd.BaseResource.Status = router.Started
	return nil
}

func (sd *StaticDirectory) Stop() error {
	sd.BaseResource.Status = router.Stopped
	return nil
}

/*
Adds handlers serving each of the static files in the whitelist under this directory. Uses the path as the base path.
*/
func (sd *StaticDirectory) AddHandlers(path string, router *httprouter.Router) error {
	for _, file := range sd.whitelist {
		var reqpath string
		if strings.HasSuffix(file, "*") {
			// httprouter requires wildcards to be named
			reqpath = path + "/" + file + "filepath"
		} else {
			reqpath = path + "/" + file
		}

		router.GET(reqpath, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			filepath := sd.path + "/" + req.URL.Path[len(path)+1:]

			// Check if the file exists and we're allowed to serve it
			info, err := os.Stat(filepath)
			if err != nil || (!sd.allow_directory_browsing && info.IsDir()) {
				http.NotFound(w, req)
				return
			}

			http.ServeFile(w, req, filepath)
		})
	}
	return nil
}
