# Makefile for WAF Tester single binary build

.PHONY: all build install clean run dev frontend backend build-all help

# Default target
all: build

# Install frontend dependencies
install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build single binary with embedded frontend and payloads
build: install
	@echo "Building WAF Tester..."
	@bash build.sh

# Build frontend only
frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Build backend only (requires frontend to be built first)
backend:
	@echo "Building backend..."
	cd backend && go build -o ../waf-tester

# Cross-compile for all platforms
build-all: install frontend
	@echo "Cross-compiling for all platforms..."
	@mkdir -p dist
	@echo "Copying frontend and payloads to backend..."
	@mkdir -p backend/frontend backend/payloads
	@cp -r frontend/dist backend/frontend/
	@cp -r payloads/*.txt backend/payloads/
	@echo "Building for macOS Intel..."
	cd backend && GOOS=darwin GOARCH=amd64 go build -o ../dist/waf-tester-darwin-amd64
	@echo "Building for macOS Apple Silicon..."
	cd backend && GOOS=darwin GOARCH=arm64 go build -o ../dist/waf-tester-darwin-arm64
	@echo "Building for Linux..."
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../dist/waf-tester-linux-amd64
	@echo "Building for Windows..."
	cd backend && GOOS=windows GOARCH=amd64 go build -o ../dist/waf-tester-windows-amd64.exe
	@echo "✓ All binaries built in ./dist/"

# Development mode - run separate servers
dev:
	@echo "Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@cd backend && go run . & cd frontend && npm run dev

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/frontend/
	rm -rf backend/payloads/
	rm -f waf-tester waf-tester.exe
	rm -rf frontend/dist/
	rm -rf dist/

# Run the built binary
run: build
	@echo "Starting WAF Tester..."
	./waf-tester

# Show help
help:
	@echo "WAF Tester - Build Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  make              Build single binary (default)"
	@echo "  make install      Install dependencies"
	@echo "  make build        Build single binary"
	@echo "  make frontend     Build frontend only"
	@echo "  make backend      Build backend only"
	@echo "  make build-all    Cross-compile for all platforms"
	@echo "  make dev          Run development servers"
	@echo "  make clean        Remove build artifacts"
	@echo "  make run          Build and run"
	@echo "  make help         Show this help"
