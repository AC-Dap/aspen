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

### Services

Sometimes, resources may be dependent on external sources; for example, an app that is still under development and changing. Rather than requiring the user to manually update these resources, Aspen can automatically fetch and deploy these resources. This requires:
- The resource is hosted on a remote Git repository
- The resource can be deployed using Docker

Aspen then deploys these servies as follows:
1. Aspen pulls the latest version of the resource from the remote Git repository into a folder `/services/<resource_id>`.
2. Aspen builds a Docker image from the Dockerfile in the resource folder.
4. Aspen starts the Docker container and maps any ports specified in the resource configuration.
  - The local port is **randomly assigned** to avoid conflicts.
5. Any existing resource can point to this service, and will be routed to the correct port.

In the code the lifecycle of a service is as follows:
1. A `Service` struct is created by parsing the config file.
2. Each service is built; this may or may not do anything, depending on whether the service has been built before.
3. Each service is started, which runs the Docker container and maps the ports.
  - If the service is already running, this does nothing.
4. The running services are passed when creating a new router instance. This lets resources create handlers using the service ports and volumns.
5. When the router is stopped, all running services are stopped and removed.

One caveat is that a service may be referenced by an old and new router. To avoid stopping anything that is still in use, we keep track of the number of routers that are using a service. When a router is stopped, it decrements the count. Only once the count reaches zero will the service be stopped and removed.
