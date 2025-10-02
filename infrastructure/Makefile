# Infrastructure Makefile
# Basic Docker orchestration commands

SHELL := /bin/bash
.SHELLFLAGS := -eu -c

DOCKER_DIR := docker

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show infrastructure commands
	@echo "Infrastructure Commands"
	@echo "Usage: make <target>"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "For SAM and Terraform commands, see: infrastructure/README.md"

.PHONY: docker-up docker-down docker-logs docker-clean
docker-up: ## Start Docker services and wait for health checks
	@echo "Starting Docker services..."
	@cd $(DOCKER_DIR) && docker compose up -d --wait
	@echo ""
	@echo "✓ Services started and healthy!"
	@echo ""
	@echo "  • API:      http://localhost:8080/api/v1"
	@echo "  • Swagger:  http://localhost:8080/swagger/index.html"
	@echo "  • DB Admin: http://localhost:8081"
	@echo ""

docker-down: ## Stop Docker services (keeps data)
	@cd $(DOCKER_DIR) && docker compose down

docker-logs: ## View Docker logs (use: make docker-logs service=app)
	@cd $(DOCKER_DIR) && docker compose logs -f $(service)

docker-clean: ## Remove all data (WARNING: deletes database!)
	@echo "⚠️  WARNING: This will delete all Docker data including the database!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		cd $(DOCKER_DIR) && docker compose down -v; \
		echo "All data removed"; \
	else \
		echo "Cancelled"; \
	fi

