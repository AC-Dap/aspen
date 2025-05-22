package router

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
	GetID() string
	GetStatus() Status
	Start() error
	Stop() error

	/*
		AddHandlers adds the resource's handlers to the router, under the given path.
	*/
	AddHandlers(path string, router *httprouter.Router) error
}

type BaseResource struct {
	Id     string
	Status Status
}

func (r *BaseResource) GetID() string {
	return r.Id
}

func (r *BaseResource) GetStatus() Status {
	return r.Status
}
