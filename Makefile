.PHONY: help build clean test install uninstall release build-all dev lint fmt

# Variables
BINARY_NAME=gomicrogen
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for current platform
	@echo "Building ${BINARY_NAME} for $(shell go env GOOS)/$(shell go env GOARCH)..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}-*
	rm -rf dist/
	rm -rf release/

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

install: build ## Build and install to /usr/local/bin
	@echo "Installing ${BINARY_NAME} to /usr/local/bin..."
	sudo cp ${BINARY_NAME} /usr/local/bin/
	sudo chmod +x /usr/local/bin/${BINARY_NAME}
	@echo "Installation complete!"

uninstall: ## Remove from /usr/local/bin
	@echo "Uninstalling ${BINARY_NAME}..."
	sudo rm -f /usr/local/bin/${BINARY_NAME}
	@echo "Uninstallation complete!"

dev: ## Run in development mode with hot reload
	@echo "Starting development server..."
	go run .

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Cross-platform builds
build-linux: ## Build for Linux (amd64, arm64)
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	chmod +x dist/${BINARY_NAME}-linux-*

build-darwin: ## Build for macOS (amd64, arm64)
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	chmod +x dist/${BINARY_NAME}-darwin-*

build-windows: ## Build for Windows (amd64, arm64)
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-arm64.exe .

build-all: clean ## Build for all platforms
	@echo "Building for all platforms..."
	mkdir -p dist
	$(MAKE) build-linux
	$(MAKE) build-darwin
	$(MAKE) build-windows
	@echo "All builds complete!"

# Release preparation
release: build-all ## Prepare release artifacts
	@echo "Preparing release artifacts..."
	mkdir -p release
	cd dist && tar -czf ../release/${BINARY_NAME}-linux-amd64.tar.gz ${BINARY_NAME}-linux-amd64
	cd dist && tar -czf ../release/${BINARY_NAME}-linux-arm64.tar.gz ${BINARY_NAME}-linux-arm64
	cd dist && tar -czf ../release/${BINARY_NAME}-darwin-amd64.tar.gz ${BINARY_NAME}-darwin-amd64
	cd dist && tar -czf ../release/${BINARY_NAME}-darwin-arm64.tar.gz ${BINARY_NAME}-darwin-arm64
	cd dist && zip ../release/${BINARY_NAME}-windows-amd64.zip ${BINARY_NAME}-windows-amd64.exe
	cd dist && zip ../release/${BINARY_NAME}-windows-arm64.zip ${BINARY_NAME}-windows-arm64.exe
	@echo "Release artifacts created in release/ directory"

# Docker builds
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -it ${BINARY_NAME}:latest

# Version info
version: ## Show version information
	@echo "Version: ${VERSION}"
	@echo "Build Time: ${BUILD_TIME}"
	@echo "Go Version: $(shell go version)" 