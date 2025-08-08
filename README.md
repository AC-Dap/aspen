# Aspen

Aspen is a highly scalable and flexible reverse proxy server that provides a unified interface for accessing multiple backend services. It's designed to be easy to configure and manage through a single JSON configuration file.

This project is still in active development, and we welcome contributions and feedback!

## Features

- **Unified Interface**: Single entrypoint that routes to multiple backend services
- **Multiple Resource Types**: Support for static files, directories, proxies, redirects, and API endpoints
- **Dynamic Service Management**: Automatically pull, build, and deploy services from Git repositories using Docker
- **Middleware Support**: Extensible middleware system for authentication, logging, and request processing
- **Hot Reload**: Update configuration without restarting the server
- **Built-in Authentication**: User management with roles and permissions (planned)

## Quick Start

### Prerequisites

- Go 1.24.3 or later
- Docker and Docker Compose (for service management)
- Git (for pulling remote services)

### Installation

```bash
git clone <repository-url>
cd aspen
go mod download
```

### Basic Usage

1. Create a configuration file (see [Configuration](#configuration) below)
2. Run Aspen:

```bash
go run aspen.go [flags] <config-file>
```

**Flags:**
- `-port`: Server port (default: 8080)
- `-services`: Folder for service files (default: ./services)

**Example:**
```bash
go run aspen.go -port 3000 -services ./my-services config.json
```

## Configuration

Aspen uses a single JSON configuration file as the source of truth. The configuration supports hot reloading, allowing you to update the server behavior without restarts.

### Configuration Structure

```json
{
  "LastUpdated": 1234567890,
  "Middleware": ["logger"],
  "Routes": [
    {
      "Route": "/api",
      "Id": "api-endpoint",
      "Resource": {
        "ResourceType": "proxy",
        "Params": {
          "Host": "http://localhost:3000",
          "Methods": ["GET", "POST"],
          "Path": "/api/*path"
        }
      }
    }
  ],
  "Services": [
    {
      "Id": "my-service",
      "Remote": "https://github.com/user/repo.git",
      "CommitHash": "abc123..."
    }
  ]
}
```

### Resource Types

Aspen supports several resource types, each handling requests differently:

#### Static File
Serves a single static file from the filesystem.

```json
{
  "ResourceType": "static_file",
  "Params": {
    "Filepath": "path/to/file.html"
  }
}
```

#### Static Directory
Serves files from a directory with optional directory browsing and file whitelisting.

```json
{
  "ResourceType": "directory",
  "Params": {
    "Path": "path/to/directory",
    "AllowDirectoryBrowsing": true,
    "Whitelist": ["*.html", "*.css", "*.js"]
  }
}
```

#### Proxy
Forwards requests to another HTTP server (local or remote).

```json
{
  "ResourceType": "proxy",
  "Params": {
    "Host": "http://localhost:8080",
    "Methods": ["GET", "POST", "PUT"],
    "Path": "/api/*path"
  }
}
```

#### Redirect
Redirects clients to another URL.

```json
{
  "ResourceType": "redirect",
  "Params": {
    "Host": "https://example.com",
    "Path": "/*path"
  }
}
```

#### API
Exposes Aspen's management API endpoints for configuration updates.

```json
{
  "ResourceType": "api",
  "Params": {}
}
```

### Services

Services are external applications that Aspen can automatically manage. They must be hosted in Git repositories and deployable with Docker.

```json
{
  "Id": "my-app",
  "Remote": "https://github.com/user/my-app.git",
  "CommitHash": "specific-commit-hash"
}
```

**Service Lifecycle:**
1. Pull source code from Git repository to `/services/<service-id>/`
2. Build Docker image from the repository's Dockerfile
3. Start container using `docker compose up -d`
4. Route requests from proxy resources to the running service
5. Automatically rebuild and redeploy when configuration updates

### Middleware

Middleware processes all requests before they reach resource handlers. Currently supported:

- `logger`: Request logging middleware

```json
{
  "Middleware": ["logger"]
}
```

## Architecture

### Core Components

- **Router**: Central request routing and middleware execution
- **Resources**: Handle specific routes and request types
- **Services**: Manage external application lifecycle
- **Config**: Configuration parsing and validation
- **Middleware**: Request/response processing pipeline

### Request Flow

1. Incoming HTTP request
2. Middleware processing (logging, auth, etc.)
3. Route matching to appropriate resource
4. Resource-specific request handling
5. Response generation and middleware post-processing

## Development

### Adding New Resource Types

1. Create a new file in `resources/`
2. Implement the `Resource` interface
3. Register the resource in `resources/register_resources.go`

### Adding New Middleware

1. Create a new file in `middleware/`
2. Implement the `Middleware` interface
3. Register the middleware in `middleware/register_middleware.go`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
