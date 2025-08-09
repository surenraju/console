# Envoy AI Gateway Console

A modern, responsive web UI for the Envoy AI Gateway API that follows the Envoy Proxy website design patterns.

## About Envoy AI Gateway

The Envoy AI Gateway is a Kubernetes-native extension of Envoy Gateway that provides intelligent routing, authentication, schema translation, and observability for Large Language Model (LLM) APIs such as OpenAI, AWS Bedrock, Azure OpenAI, and GCP Vertex AI.

## Project Structure

```
envoy-aigateway-console/
├── frontend/                 # Vite + React app
├── backend/                  # Golang backend
├── manifests/               # Kubernetes manifests
├── docker/                  # Docker files
├── scripts/                 # Build/deploy scripts
├── docs/                    # Documentation
└── README.md
```

## Quick Start

### Development
```bash
# Run both frontend and backend
make dev

# Run tests
make test

# Build Docker images
make docker-build
```

### Docker
```bash
# Run all-in-one container
docker run -p 80:80 envoy-aigateway-console:latest
```

## Features

- **Responsive Design**: Mobile-first approach with Tailwind CSS
- **Real-time Updates**: WebSocket connection for live updates
- **Kubernetes Integration**: Direct API integration with watch capabilities
- **Modern UI**: Built with shadcn/ui components
- **Type Safety**: Full TypeScript support

## Technology Stack

- **Frontend**: Vite + React 18 + TypeScript + Tailwind CSS + shadcn/ui
- **Backend**: Golang + WebSocket
- **Containerization**: Docker + Caddy
- **Deployment**: Kubernetes + Helm

## Development

See [docs/development.md](docs/development.md) for detailed development setup instructions. 