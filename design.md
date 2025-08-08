### Design Goals

Aspen is a reverse proxy server that provides a unified interface for accessing multiple backend services. It is designed to be highly scalable, flexible, and easy to use.

We have a single entrypoint that then branches into many **resources** hosted on some **route**. Resources can be static files, local servers, or redirects to external URLs. Each resource is configured with a set of rules that determine how requests are routed to the backend services.

Additionally, Aspen manages many **services** that can be dynamically started and stopped. These services are typically backend applications that are hosted in Docker containers. Aspen can automatically pull the latest version of these services from a remote Git repository, build them, and start them as needed.

Aspen also contains a built in authentication system that can be used to secure access to resources. This system supports users with passwords and different roles/permissions.

All of this configuration lies in a **single JSON file** that is read at startup. Aspen is designed to treat this file as the **source of truth** for the configuration of the server, and any changes to the configuration should be made in this file. Additionally, Aspen should be able to reload the configuration file without restarting the server, allowing for dynamic updates to the server's behavior.

### Routes and Resources
Resources are the main components of Aspen. They represent the various backend services that are hosted on the server. There are a few different resource **types** that can be configured with Aspen, and are added to the server under a specific route. This resource is then **completely responsible** for handling requests that come to that route.

Internally, resources correspond to a set of http handlers that are registered with the router. Resources are free to add any number of handlers, and these handlers can be used to process requests in any way the resource needs. For example, a resource may have a handler that serves static files, or many handlers that process a variety of different API requests.

Currently, resources can be of the following types:
- **Static file**: This is an existing file on the server that is directly served to the client.
- **Static directory**: This is an existing directory on the server that is directly served to the client. The resource serves a whitelist of files to the client, and optionally allows directory browsing.
- **Proxy**: This proxies requests to another URL. This is primarily used to host services that are locally running, but can also be used to proxy requests to external URLs.
- **Redirect**: This redirects the client to another URL.
- **API**: This surfaces API endpoints that can be used to configure Aspen. They can allow clients to modify the configuration file and reload Aspen to apply the changes. This is primarily used for the admin dashboard, but can also be used to expose other API endpoints.

### Services

Sometimes, resources may be dependent on external sources; for example, an app that is still under development and changing. Rather than requiring the user to manually update these resources, Aspen can automatically fetch and deploy these resources. This requires:
- The resource is hosted on a remote Git repository
- The resource can be deployed using Docker

Aspen then deploys these servies as follows:
1. Aspen pulls the latest version of the resource from the remote Git repository into a folder `/services/<resource_id>`.
2. Aspen builds a Docker image from the Dockerfile in the resource folder.
4. Aspen starts the Docker container with `docker compose`.
5. Any existing resource can point to this service, and will be routed to the correct port.

In the code the lifecycle of a service is as follows:
1. A `Service` struct is created by parsing the config file.
2. The services are passed when creating a new router instance. This lets resources reference the service folder.
  - Note that services are not running yet!
3. Each service is built; if the current version has already been built, this does nothing.
  - If the service is not built, it is pulled from the remote Git repository and built into a Docker image.
  - If the most recently built version is out of date, it is pulled and rebuilt.
4. Each service is started using `docker compose up -d`.
  - If the service is already running, this does nothing
  - If an out of date service is running, Docker will automatically stop the old container and start a new one with the latest version.
5. When the router is stopped, all running services are stopped and removed.

One caveat is that a service may be referenced by an old and new router. To avoid stopping anything that is still in use, we keep track of the number of routers that are using a service. When a router is stopped, it decrements the count. Only once the count reaches zero will the service be stopped and removed.

## Middleware

Aspen supports middleware that can be used to process requests before they reach the resource handlers. Middleware can be used to perform tasks such as authentication, logging, and request modification.

Middleware is applied to **every request** on **every route**. This means that middleware should try and be as efficient as possible, and should not perform any blocking operations. Middleware can be used to modify the request or response, or to perform any other necessary processing.
