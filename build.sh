#!/bin/bash

# Build script for Confused Dependency Confusion Scanner

echo "Building Confused Dependency Confusion Scanner..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Error: Go version $GO_VERSION is not supported. Please install Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "Go version: $GO_VERSION"

# Clean previous builds
echo "Cleaning previous builds..."
rm -f confused confused.exe

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Build for current platform
echo "Building for current platform..."
go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused ./cmd/confused

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: confused"
    echo "You can now run: ./confused --help"
else
    echo "Build failed!"
    exit 1
fi

# Optional: Build for multiple platforms
if [ "$1" = "--cross-compile" ]; then
    echo "Building for multiple platforms..."
    
    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused-linux-amd64 ./cmd/confused
    
    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused-windows-amd64.exe ./cmd/confused
    
    # macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused-darwin-amd64 ./cmd/confused
    
    # macOS ARM64
    GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused-darwin-arm64 ./cmd/confused
    
    echo "Cross-compilation complete!"
fi

echo "Build process completed!"
