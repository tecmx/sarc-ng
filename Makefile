SHELL := /bin/bash
.SHELLFLAGS := -eu -c
.DEFAULT_GOAL := help

# Variables
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
APP_BINARY := $(BIN_DIR)/app
CLI_BINARY := $(BIN_DIR)/sarc

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildDate=$(BUILD_DATE)'

APP_MAIN := ./cmd/server
CLI_MAIN := ./cmd/cli

# Helper functions
define check_tool
	@command -v $(1) >/dev/null || (echo "Error: $(1) not installed. Run 'make setup'" && exit 1)
endef

# Helper target to conditionally run generate
.PHONY: _maybe_generate
_maybe_generate:
ifndef SKIP_GENERATE
	@$(MAKE) generate
endif

# Main targets
.PHONY: help setup
help: ## Show available commands
	@echo "SARC-NG Development Makefile"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup development environment and install dependencies
	@go mod download && go mod tidy
	@grep -E '^\s*_\s+"[^"]+"\s*//.*' tools.go | sed 's/.*"\([^"]*\)".*/\1/' | xargs -I {} go install {}@latest

# Development
.PHONY: run debug wire generate
run: _maybe_generate ## Run the application directly
	go run $(APP_MAIN)

debug: _maybe_generate ## Run with hot reloading (requires air)
	$(call check_tool,air)
	air -c .air.toml

wire: ## Generate dependency injection code
	@go generate ./cmd/server ./cmd/lambda
	@test -f cmd/server/wire_gen.go && test -f cmd/lambda/wire_gen.go || (echo "Wire generation failed" && exit 1)

generate: wire swagger ## Generate dependency injection code and Swagger docs

swagger: ## Generate API documentation
	$(call check_tool,swag)
	@rm -rf api/swagger/swagger.json api/swagger/swagger.yaml api/swagger/docs.go
	swag init -g cmd/server/main.go --parseDependency --parseInternal --output api/swagger
	@test -f api/swagger/swagger.json || (echo "Swagger generation failed" && exit 1)

# Build
.PHONY: build
build: _maybe_generate bin-dir ## Build server and CLI applications
	go build -ldflags="$(LDFLAGS)" -o $(APP_BINARY) $(APP_MAIN)
	go build -ldflags="$(LDFLAGS)" -o $(CLI_BINARY) $(CLI_MAIN)

# Quality Assurance
.PHONY: format lint test
format: ## Format Go code
	go fmt ./...

lint: _maybe_generate ## Run linters
	go vet ./...
	$(call check_tool,golangci-lint)
	golangci-lint run

test: _maybe_generate ## Run tests
	go test -race ./...

# Cleanup
.PHONY: clean
clean: ## Remove all build artifacts, caches, and generated files
	rm -rf $(BUILD_DIR) dist tmp
	rm -f coverage.out coverage.html
	rm -f cmd/server/wire_gen.go cmd/lambda/wire_gen.go
	rm -f api/swagger/docs.go api/swagger/swagger.json api/swagger/swagger.yaml
	rm -rf docs/.docusaurus docs/build
	go clean -cache

# Internal helpers
.PHONY: bin-dir
bin-dir:
	@mkdir -p $(BIN_DIR)
