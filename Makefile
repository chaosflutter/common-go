.PHONY: help dev build run clean test install-air

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

dev: ## Start development server with Air (live reload)
	@./dev.sh

build: ## Build the application
	@echo "Building application..."
	@go build -o app main.go
	@echo "✅ Build complete! Binary saved as 'app'"

run: build ## Build and run the application
	@echo "🚀 Starting server..."
	@./app

clean: ## Clean build artifacts and temporary files
	@echo "🧹 Cleaning up..."
	@rm -f app
	@rm -rf tmp/
	@rm -f build-errors.log
	@echo "✅ Clean complete!"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

install-air: ## Install Air for live reload
	@echo "📦 Installing Air..."
	@go install github.com/air-verse/air@latest
	@echo "✅ Air installed successfully!"

deps: ## Download and tidy up dependencies
	@echo "📦 Downloading dependencies..."
	@go mod tidy
	@echo "✅ Dependencies updated!"

fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete!"

check: fmt vet test ## Run all checks (format, vet, test)
	@echo "✅ All checks passed!"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t common-go-app .
	@echo "✅ Docker image built!"

docker-run: ## Run application in Docker
	@echo "🐳 Running application in Docker..."
	@docker run -p 8080:8080 \
		-v $(shell pwd)/logs:/root/logs \
		-v $(shell pwd)/audio:/root/audio \
		-v $(shell pwd)/.env:/root/.env:ro \
		common-go-app

docker-run-detached: ## Run application in Docker (detached mode)
	@echo "🐳 Running application in Docker (detached)..."
	@docker run -d -p 8080:8080 \
		-v $(shell pwd)/logs:/root/logs \
		-v $(shell pwd)/audio:/root/audio \
		-v $(shell pwd)/.env:/root/.env:ro \
		--name common-go-container common-go-app
	@echo "✅ Container started! Access at http://localhost:8080"

docker-stop: ## Stop Docker container
	@echo "🐳 Stopping Docker container..."
	@docker stop common-go-container || true
	@docker rm common-go-container || true
	@echo "✅ Container stopped!"

docker-logs: ## View Docker container logs
	@docker logs -f common-go-container

docker-compose-up: ## Run with docker-compose
	@echo "🐳 Starting with docker-compose..."
	@docker-compose up --build

docker-compose-up-detached: ## Run with docker-compose (detached)
	@echo "🐳 Starting with docker-compose (detached)..."
	@docker-compose up -d --build
	@echo "✅ Services started! Access at http://localhost:8080"

docker-compose-down: ## Stop docker-compose services
	@echo "🐳 Stopping docker-compose services..."
	@docker-compose down
	@echo "✅ Services stopped!"

docker-clean: ## Clean Docker images and containers
	@echo "🐳 Cleaning Docker resources..."
	@docker stop common-go-container || true
	@docker rm common-go-container || true
	@docker rmi common-go-app || true
	@docker-compose down || true
	@echo "✅ Docker resources cleaned!"

test-api: ## Test all API endpoints
	@echo "🧪 Testing API endpoints..."
	@./test_api.sh

tts-examples: ## Run TTS API examples
	@echo "🎤 Running TTS API examples..."
	@./examples/tts_example.sh

setup: deps install-air ## Initial project setup
	@echo "🎉 Project setup complete! Run 'make dev' to start developing."
