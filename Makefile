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
	@echo "📱 Frontend: http://localhost:5173"
	@echo "🔧 Backend: http://localhost:8080"
	@make -j2 frontend backend

# Frontend development
frontend:
	@echo "🎨 Starting frontend development server..."
	cd frontend && yarn dev

# Backend development
backend:
	@echo "⚙️  Starting backend development server..."
	cd backend && go run cmd/server/main.go

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
	@echo "  make frontend   - Start frontend development server"
	@echo "  make backend    - Start backend development server"
	@echo "  make build      - Build both frontend and backend"
	@echo "  make test       - Run all tests"
	@echo "  make install    - Install all dependencies"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up  - Start Docker services"
	@echo "  make docker-down - Stop Docker services" 