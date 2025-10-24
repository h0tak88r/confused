# Multi-stage Dockerfile for Confused Dependency Confusion Scanner

# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -ldflags "-X main.version=2.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o confused ./cmd/confused

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    curl \
    git \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh confused

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/confused /usr/local/bin/confused

# Copy configuration files
COPY confused.yaml /etc/confused/confused.yaml

# Create results directory
RUN mkdir -p /app/results && chown confused:confused /app/results

# Switch to non-root user
USER confused

# Set environment variables
ENV CONFUSED_OUTPUT_DIR=/app/results
ENV CONFUSED_VERBOSE=false

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD confused --help > /dev/null || exit 1

# Default command
ENTRYPOINT ["confused"]
CMD ["--help"]

# Labels
LABEL maintainer="h0tak88r"
LABEL version="2.0.0"
LABEL description="Advanced Dependency Confusion Scanner"
LABEL org.opencontainers.image.title="Confused"
LABEL org.opencontainers.image.description="Advanced Dependency Confusion Scanner"
LABEL org.opencontainers.image.version="2.0.0"
LABEL org.opencontainers.image.source="https://github.com/h0tak88r/confused"
