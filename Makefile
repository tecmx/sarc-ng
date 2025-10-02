# Define shell and options - use portable shell detection
SHELL := /bin/bash
.SHELLFLAGS := -eu -c

#
# SARC-NG Development Makefile
# ============================
# Go-based API for resource management system
#
# Quick Start:
#   make setup        - Install dependencies and tools
#   make build        - Build server and CLI applications
#   make help         - Show all available commands
#
# For Docker deployment:
#   cd infrastructure/docker && docker compose up -d
#

#

# Variables
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
APP_BINARY := $(BIN_DIR)/app
CLI_BINARY := $(BIN_DIR)/sarc
COVERAGE_OUT := $(BUILD_DIR)/coverage.out
COVERAGE_HTML := $(BUILD_DIR)/coverage.html

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build configuration
LDFLAGS := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.buildDate=$(BUILD_DATE)'

APP_MAIN := ./cmd/server
CLI_MAIN := ./cmd/cli

# Default target
.DEFAULT_GOAL := help

# Helper functions
define check_tool
	@command -v $(1) >/dev/null || (echo "Error: $(1) not installed. Run 'make setup'" && exit 1)
endef

#
# MAIN TARGETS
#

.PHONY: help setup
help: ## Show available commands
	@echo "SARC-NG Development Makefile"
	@echo "Usage: make <target>"
	@echo ""
	@echo "Development Commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-up       Start Docker services"
	@echo "  docker-down     Stop Docker services"
	@echo "  docker-logs     View logs (add service=<name> for specific service)"
	@echo "  docker-clean    Remove all data (WARNING: deletes database)"
	@echo ""
	@echo "For SAM and Terraform, see: infrastructure/README.md"

setup: ## Setup development environment and install dependencies
	@echo "Setting up development environment..."
	@go mod download && go mod tidy
	@echo "Installing development tools..."
	@grep -E '^\s*_\s+"[^"]+"\s*//.*' tools.go | sed 's/.*"\([^"]*\)".*/\1/' | xargs -I {} go install {}@latest
	@echo "Setup completed!"

#
# DEVELOPMENT
#

.PHONY: run debug wire
run: ## Run the application directly
	go run $(APP_MAIN)

debug: ## Run with hot reloading (requires air)
	$(call check_tool,air)
	air -c .air.toml

wire: ## Generate dependency injection code
	@echo "Generating Wire dependency injection code..."
	@go generate ./cmd/server
	@go generate ./cmd/lambda
	@echo "Verifying generated files..."
	@test -f cmd/server/wire_gen.go || (echo "✗ Failed: cmd/server/wire_gen.go not generated" && exit 1)
	@test -f cmd/lambda/wire_gen.go || (echo "✗ Failed: cmd/lambda/wire_gen.go not generated" && exit 1)
	@echo "✓ Wire code generated successfully"



#
# BUILD & RELEASE
#

.PHONY: build release
build: wire bin-dir ## Build server and CLI applications
	@echo "Building server application..."
	go build -ldflags="$(LDFLAGS)" -o $(APP_BINARY) $(APP_MAIN)
	@echo "Server built: $(APP_BINARY)"
	@echo "Building CLI application..."
	go build -ldflags="$(LDFLAGS)" -o $(CLI_BINARY) $(CLI_MAIN)
	@echo "CLI built: $(CLI_BINARY)"

release: ## Build production release binaries
	@echo "Building production release..."
	@mkdir -p $(BIN_DIR)
	@echo "Building server for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/server-linux-amd64 $(APP_MAIN)
	@echo "Building CLI for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/sarc-linux-amd64 $(CLI_MAIN)
	@echo "Production binaries built in $(BIN_DIR)"

#
# DOCUMENTATION
#

.PHONY: swagger docs-build docs-serve docs-clean
swagger: ## Generate API documentation
	@echo "Generating Swagger docs..."
	$(call check_tool,swag)
	@rm -rf api/swagger/swagger.json api/swagger/swagger.yaml api/swagger/docs.go
	swag init -g cmd/server/main.go --parseDependency --parseInternal --output api/swagger
	@[ -f api/swagger/swagger.json ] || (echo "Generation failed" && exit 1)
	@grep -q '"paths": {}' api/swagger/swagger.json && echo "Empty paths - check annotations" || echo "Documentation generated"

docs-build: ## Build documentation site
	@echo "Building documentation..."
	$(call check_tool,yarn)
	@cd docs && yarn clean-api-docs sarc-ng || true
	@cd docs && yarn gen-api-docs sarc-ng
	@cd docs && yarn build
	@echo "Documentation built successfully in docs/build/"

docs-serve: ## Start documentation development server
	@echo "Starting documentation server..."
	$(call check_tool,yarn)
	@cd docs && yarn start

docs-clean: ## Clean documentation build artifacts
	@echo "Cleaning documentation..."
	@cd docs && yarn clean-api-docs sarc-ng || true
	@rm -rf docs/build docs/node_modules/.cache
	@echo "Documentation cleaned"

#
# QUALITY ASSURANCE
#

.PHONY: format lint test coverage
format: ## Format Go code
	go fmt ./...

lint: ## Run linters
	go vet ./...
	$(call check_tool,golangci-lint)
	golangci-lint run

test: ## Run tests
	go test -race ./...

coverage: ## Generate test coverage report
	go test -race -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"

#
# WORKFLOW HELPERS
#

.PHONY: pre-commit check ci
pre-commit: format wire lint test ## Run all pre-commit checks
	@echo "✓ All pre-commit checks passed"

check: wire lint test ## Quick validation (format, lint, test)
	@echo "✓ Code checks passed"

ci: ## CI pipeline (wire, lint, test, build)
	@echo "Running CI pipeline..."
	@$(MAKE) wire
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) build
	@echo "✓ CI pipeline completed"

#
# CLEANUP
#

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) dist .env tmp
	go clean -cache
	@echo "Cleanup completed"

.PHONY: clean-all
clean-all: clean ## Clean all artifacts (codebase only)
	@echo "All codebase artifacts cleaned"
	@echo "For Docker cleanup: make docker-clean"

#
# INFRASTRUCTURE DELEGATION
#

.PHONY: docker-up docker-down docker-logs docker-clean
docker-up docker-down docker-logs docker-clean:
	@$(MAKE) -C infrastructure $@

# Internal helpers
.PHONY: bin-dir
bin-dir:
	@mkdir -p $(BIN_DIR)
