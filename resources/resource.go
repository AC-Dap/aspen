package resources

import (
	"github.com/julienschmidt/httprouter"
)

type Status int

const (
	NotStarted Status = iota
	Started
	Stopping
	Stopped
)

type Resource interface {
	Start() error
	Stop() error
	Status() Status

	/*
		AddHandlers adds the resource's handlers to the router, under the given path.
	*/
	AddHandlers(path string, router *httprouter.Router) error
}
