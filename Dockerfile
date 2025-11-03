# Multi-stage Dockerfile for production deployment
# Optimized for Kubernetes environments

# Stage 1: Builder
FROM golang:1.24.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary (no C dependencies - important for scratch/alpine)
# -ldflags for optimization and version info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -a -installsuffix cgo \
    -o citary-backend \
    ./cmd/api

# Stage 2: Final lightweight image
FROM alpine:3.19

# Add metadata labels (useful for Kubernetes)
LABEL maintainer="your-email@example.com"
LABEL app="citary-backend"
LABEL version="1.0.0"

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && update-ca-certificates

# Create non-root user for security (Kubernetes best practice)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/citary-backend .

# Copy timezone data (if your app needs specific timezones)
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user (Kubernetes security requirement)
USER appuser

# Expose port (should match your app's PORT env var)
EXPOSE 3001

# Health check (Kubernetes will use this for liveness/readiness probes)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3001/health || exit 1

# Set environment variables (can be overridden by Kubernetes)
ENV PORT=3001

# Run the application
ENTRYPOINT ["/app/citary-backend"]
