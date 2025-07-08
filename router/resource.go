package router

type Resource interface {
	GetID() string

	/*
		AddHandlers adds the resource's handlers to the router, under the given path.
	*/
	AddHandlers(path string, router *RouterInstance) error
}

type BaseResource struct {
	id string
}

// NewBaseResource creates a new BaseResource with the given ID.
func NewBaseResource(id string) BaseResource {
	return BaseResource{
		id: id,
	}
}

func (r *BaseResource) GetID() string {
	return r.id
}
