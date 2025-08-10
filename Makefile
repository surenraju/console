.PHONY: clean dev build test frontend backend

# Clean all build artifacts and dependencies
clean:
	@echo "🧹 Cleaning frontend..."
	cd frontend && yarn clean
	@echo "🧹 Cleaning backend..."
	cd backend && go clean -cache -testcache -modcache
	@echo "✅ Clean complete!"

# Development commands
dev: clean
	@echo "🚀 Starting development environment..."
	@echo "📱 Frontend: http://localhost:5173 (or next available port)"
	@echo "🔧 Backend: http://localhost:8082"
	@echo "📋 Use 'make logs' to see detailed logs"
	@echo "📋 Use 'make stop' to stop all services"
	@make -j2 frontend backend

# Start development servers with logging
dev-logs: clean
	@echo "🚀 Starting development environment with detailed logs..."
	@make -j2 frontend-verbose backend-verbose

# Frontend development
frontend:
	@echo "🎨 Starting frontend development server..."
	cd frontend && yarn dev

# Frontend development with verbose output  
frontend-verbose:
	@echo "🎨 Starting frontend development server (verbose)..."
	cd frontend && yarn dev --host 0.0.0.0

# Backend development
backend:
	@echo "⚙️  Starting backend development server..."
	cd backend && SERVER_PORT=8082 go run cmd/server/main.go

# Backend development with verbose output
backend-verbose:
	@echo "⚙️  Starting backend development server (verbose)..."
	cd backend && SERVER_PORT=8082 go run cmd/server/main.go -v

# Check if services are running
status:
	@echo "📊 Checking service status..."
	@echo "Frontend (port 5173/5174):"
	@lsof -i :5173 -i :5174 | grep LISTEN || echo "  ❌ Frontend not running"
	@echo "Backend (port 8082):"
	@lsof -i :8082 | grep LISTEN || echo "  ❌ Backend not running"

# Stop all development services
stop:
	@echo "🛑 Stopping development services..."
	@pkill -f "yarn dev" || true
	@pkill -f "vite" || true  
	@pkill -f "go run cmd/server/main.go" || true
	@echo "✅ Services stopped"

# Quick start (skip clean)
quick:
	@echo "⚡ Quick starting development environment..."
	@echo "📱 Frontend: http://localhost:5173 (or next available port)"
	@echo "🔧 Backend: http://localhost:8082"
	@make -j2 frontend backend

# Build both frontend and backend
build:
	@echo "🔨 Building frontend..."
	cd frontend && yarn build
	@echo "🔨 Building backend..."
	cd backend && go build -o bin/server cmd/server/main.go
	@echo "✅ Build complete!"

# Run tests
test:
	@echo "🧪 Running frontend tests..."
	cd frontend && yarn test
	@echo "🧪 Running backend tests..."
	cd backend && go test ./...
	@echo "✅ Tests complete!"

# Install dependencies
install:
	@echo "📦 Installing frontend dependencies..."
	cd frontend && yarn install
	@echo "📦 Installing backend dependencies..."
	cd backend && go mod download
	@echo "✅ Dependencies installed!"

# Docker commands
docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build

docker-up:
	@echo "🐳 Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker services..."
	docker-compose down

# Help
help:
	@echo "Available commands:"
	@echo "  make clean      - Clean all build artifacts"
	@echo "  make dev        - Start development environment (frontend + backend)"
	@echo "  make quick      - Quick start development (skip clean)"
	@echo "  make dev-logs   - Start development with verbose logging"
	@echo "  make frontend   - Start frontend development server only"
	@echo "  make backend    - Start backend development server only"
	@echo "  make status     - Check if services are running"
	@echo "  make stop       - Stop all development services"
	@echo "  make build      - Build both frontend and backend"
	@echo "  make test       - Run all tests"
	@echo "  make install    - Install all dependencies"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up  - Start Docker services"
	@echo "  make docker-down - Stop Docker services" 