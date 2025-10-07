# Multi-stage Dockerfile for SARC-NG
# Build development image: docker build --target development -t sarc-ng:dev .
# Build production image: docker build --target production -t sarc-ng:prod .

# Base stage with common dependencies
FROM golang:1.25-alpine AS base

# Install common system dependencies and update packages for security
RUN apk update && apk upgrade && apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Set working directory
WORKDIR /app

# Copy dependency files for better caching
COPY go.mod go.sum ./
RUN go mod download

# Development stage
FROM base AS development

# Install development tools
RUN apk add --no-cache \
    curl \
    make \
    gcc \
    libc-dev \
    shadow

# Install Go development tools (pinned versions for compatibility)
RUN go install github.com/air-verse/air@v1.52.3 && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Copy Air configuration
COPY .air.toml .air.toml

# Set CGO for debugging support
ENV CGO_ENABLED=1

# Create user with configurable UID/GID for development
ARG USER_ID=1000
ARG GROUP_ID=1000
ARG USERNAME=developer

# Create group and user with specified IDs
RUN addgroup -g ${GROUP_ID} ${USERNAME} && \
    adduser -D -u ${USER_ID} -G ${USERNAME} ${USERNAME}

# Create cache directory and ensure the user owns the app directory and Go paths
RUN mkdir -p /go/cache /go/pkg/mod && \
    chown -R ${USERNAME}:${USERNAME} /app /go

# Switch to the created user
USER ${USERNAME}

# Application code will be mounted as volume for hot reloading
CMD ["air", "-c", ".air.toml"]

# Build stage for production
FROM base AS builder

# Copy source code
COPY . .

# Generate Wire dependency injection code before building
RUN go generate ./cmd/server ./cmd/lambda

# Build the application with production optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o sarc-server ./cmd/server

# Production stage
FROM gcr.io/distroless/static:nonroot AS production

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/sarc-server .

# Use non-root user for security
USER nonroot:nonroot

# Run the application
CMD ["./sarc-server"]
