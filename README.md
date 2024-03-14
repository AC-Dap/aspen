### Config Files

#### `server.toml`
This file contains any necessary global configuration information for the server.
This includes:
- The server's port
- The nginx port
- The 404 page

#### `resources.toml`
This contains all the hosted resources.
Each resource has:
- A unique name
- A unique route
- A resource definition
  - Types:
    - Static, self-hosted. This is an existing file on the server (e.g. dashboard).
    - Static, external. This is a file that should be pulled from an external source (e.g. github).
    - Dynamic, self-hosted. This is a localhost server that is externally started
    - Dynamic, external. This is a server that should be pulled from an external source (e.g. github).
  - Path/port
  - External source (if applicable)
  - Command to build + start
    - npm, docker, bash script, etc.
  - Command to stop
    - npm, docker, bash script, etc.
  - Command to update
    - npm, docker, bash script, etc.
- Whether authorization is needed

Each entry should be a single table under the `resources` table array, like so:
```toml
[[resources]]
name = "dashboard"
route = "/dashboard"
...

[[resources]]
name = "auth"
route = "/auth"
...
```
**We require the existence of a `dashboard` and `auth` resource.**