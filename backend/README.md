# Envoy AI Gateway Console - Backend

This is the backend service for the Envoy AI Gateway Console, providing REST APIs and WebSocket connectivity for managing AI Gateway configurations in Kubernetes.

## Architecture

The backend is structured as follows:

- **cmd/server/** - Main application entry point
- **internal/api/** - REST API handlers and WebSocket server
- **internal/config/** - Configuration management
- **internal/k8s/** - Kubernetes client and CRD operations
- **internal/models/** - AI Gateway CRD data models
- **internal/handlers/** - Business logic handlers
- **pkg/types/** - Shared types and interfaces

## Features

- REST API for managing AI Gateway providers and routes
- WebSocket server for real-time updates
- Kubernetes CRD integration for AI Gateway resources
- Provider templates for major LLM services (OpenAI, AWS, GCP, Azure)
- Configuration validation and error handling

## Getting Started

### Prerequisites

- Go 1.21 or later
- Access to a Kubernetes cluster with Envoy AI Gateway installed
- kubectl configured to access the cluster

### Development

1. Build the application:
   ```bash
   go build -o bin/server cmd/server/main.go
   ```

2. Run the server:
   ```bash
   ./bin/server
   ```

3. The server will start on `http://localhost:8080` by default.

### Configuration

The server can be configured using environment variables:

- `SERVER_ADDRESS` - Server bind address (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `K8S_IN_CLUSTER` - Use in-cluster Kubernetes config (default: false)
- `K8S_CONFIG_PATH` - Path to kubeconfig file (default: ~/.kube/config)
- `K8S_NAMESPACE` - Kubernetes namespace to operate in (default: default)

### API Endpoints

#### Providers
- `GET /api/v1/providers` - List all providers
- `POST /api/v1/providers` - Create a new provider
- `GET /api/v1/providers/{name}` - Get a specific provider
- `PUT /api/v1/providers/{name}` - Update a provider
- `DELETE /api/v1/providers/{name}` - Delete a provider

#### Routes
- `GET /api/v1/routes` - List all routes
- `POST /api/v1/routes` - Create a new route
- `GET /api/v1/routes/{name}` - Get a specific route
- `PUT /api/v1/routes/{name}` - Update a route
- `DELETE /api/v1/routes/{name}` - Delete a route

#### Templates
- `GET /api/v1/templates` - List available provider templates
- `GET /api/v1/templates/{provider}` - Get a specific provider template

#### WebSocket
- `GET /ws` - WebSocket endpoint for real-time updates

### Testing

Run tests with:
```bash
go test ./...
```

## Docker

Build Docker image:
```bash
docker build -t envoy-ai-gateway-console-backend .
```

Run with Docker:
```bash
docker run -p 8080:8080 envoy-ai-gateway-console-backend
```

## Development Status

This is the initial implementation following the task-based development approach. The current implementation provides:

✅ Basic project structure
✅ Configuration management  
✅ HTTP server with Gin
✅ Kubernetes client setup
✅ Basic API endpoints (stubs)
✅ WebSocket server foundation

🚧 Next steps:
- Implement AI Gateway CRD models
- Add CRD watchers and event handling
- Implement full CRUD operations
- Add validation logic
- Create comprehensive tests