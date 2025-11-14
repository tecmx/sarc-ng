# SARC-NG

[![CI](https://github.com/tecmx/sarc-ng/workflows/CI/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/ci.yml)
[![Deploy](https://github.com/tecmx/sarc-ng/workflows/Deploy/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/deploy.yml)
[![Release](https://github.com/tecmx/sarc-ng/workflows/Release/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecmx/sarc-ng)](https://goreportcard.com/report/github.com/tecmx/sarc-ng)

Resource Management and Scheduling API — a Go service for managing buildings, classes, lessons, resources, and reservations.

## Purpose

- **Centralized scheduling:** keep academic spaces, resources, and lessons in sync.
- **API-first operations:** consistent `/api/v1` contracts for every entity.
- **Cloud-ready workflows:** Docker, SAM, and Terraform stacks ship with the repo.

## Quick Start

```bash
git clone https://github.com/tecmx/sarc-ng.git
cd sarc-ng
make docker-up
```

**Access Points:**
- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html

## Structure

```
cmd/            # Entrypoints (CLI, Lambda, server)
internal/       # Domain, services, adapters, transport
infrastructure/ # Docker, SAM, Terraform
docs/           # Detailed guides
```

## Commands

```bash
make run        # Start API locally
make test       # Run unit tests
make build      # Build binaries
make docker-up  # Bring up full stack
make lint       # Static analysis
make swagger    # Refresh API docs
```

## Documentation

- [Getting Started](docs/content/getting-started.md)
- [Development Guide](docs/content/development.md)
- [Architecture](docs/content/architecture.md)
- [Deployment & CI/CD](docs/content/deployment.md#cicd-pipeline)
- [Authentication Suite](docs/content/authentication.md)
