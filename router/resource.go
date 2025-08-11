package router

type Resource interface {
	GetID() string

	/*
		AddHandlers adds the resource's handlers to the router, under the given path.
	*/
	AddHandlers(path string, router *RouterInstance) error
}

type BaseResource struct {
	Id string

	// AccessRoles is a list of roles that can access this resource.
	// If empty, the resource is accessible to all roles.
	AccessRoles []string
}

// NewBaseResource creates a new BaseResource with the given ID.
func NewBaseResource(id string, accessRoles []string) BaseResource {
	return BaseResource{
		Id:          id,
		AccessRoles: accessRoles,
	}
}

func (br BaseResource) GetID() string {
	return br.Id
}
