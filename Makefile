# Makefile for Confused Dependency Confusion Scanner

# Variables
BINARY_NAME=confused
VERSION=2.1.0
BUILD_DATE=$(shell date -u +%Y-%m-%d)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go mod tidy
	@go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/confused
	@echo "Build complete! Binary: $(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@rm -f confused-*
	@echo "Clean complete!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@go vet ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Cross-compile for multiple platforms
.PHONY: cross-compile
cross-compile: clean
	@echo "Cross-compiling for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/confused
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe ./cmd/confused
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/confused
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 ./cmd/confused
	@echo "Cross-compilation complete!"

# Run the binary
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) --help

# Generate sample configuration
.PHONY: config
config:
	@echo "Generating sample configuration..."
	@./$(BINARY_NAME) config generate -o confused.sample.yaml

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd/confused
	@echo "Installation complete!"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-race     - Run tests with race detection"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  deps          - Install dependencies"
	@echo "  cross-compile - Cross-compile for multiple platforms"
	@echo "  run           - Build and run the binary"
	@echo "  config        - Generate sample configuration"
	@echo "  install       - Install the binary to GOPATH/bin"
	@echo "  help          - Show this help message"
