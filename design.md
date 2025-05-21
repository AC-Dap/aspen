### Design

Aspen is a reverse proxy server that provides a unified interface for accessing multiple backend services. It is designed to be highly scalable, flexible, and easy to use.

We have a single entrypoint that then branches into many **resources** based on the request path. Resources can be static files, local servers, or redirects to external URLs. Each resource is configured with a set of rules that determine how requests are routed to the backend services.

Aspen also contains a built in authentication system that can be used to secure access to resources. This system supports users with passwords, and different roles/permissions. Additionally, this authentication service is available to other services within the Aspen ecosystem.

All of these features can be configured using TOML files.

### Requirements
- Configurable resources that map a request path pattern to a handler function
- An HTTP router that can handle incoming requests and route them to the appropriate resource
- A red/blue system that can swap out the router and resources as needed
- An authentication system that maps users to roles and permissions
- A logging system that logs all requests and responses
