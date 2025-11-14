# Infrastructure

Infrastructure management for SARC-NG: Docker, AWS SAM, and Terraform.

## Quick Start

```bash
make docker-up    # Start local development environment
```

## Prerequisites

- Docker - Local development
- AWS CLI - Cloud deployment (with credentials configured)
- AWS SAM CLI - Lambda deployment (`pip install aws-sam-cli`)
- Terraform - Infrastructure provisioning (optional)

## Docker (Local Development)

### Commands

```bash
make docker-up       # Start services
make docker-down     # Stop services
make docker-logs     # View logs
make docker-clean    # Remove all data
```

### Access Points

- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html
- DB Admin: http://localhost:8081
- Metrics: http://localhost:8080/metrics

📖 **[Complete Docker Guide](../docs/infrastructure/DOCKER.md)**

## AWS SAM (Serverless)

### Local Development

```bash
cd infrastructure/sam
make local-dev      # Start local API + MySQL
make local-stop     # Stop local services
```

### AWS Deployment

```bash
cd infrastructure/sam
make deploy-all ENV=dev    # Deploy to AWS
make urls ENV=dev          # Show deployed URLs
make delete-all ENV=dev    # Delete stacks
```

📖 **[Complete SAM Guide](../docs/infrastructure/SAM-GUIDE.md)**
📖 **[SAM Detailed README](sam/README.md)**

## Terraform

### Basic Workflow

```bash
cd infrastructure/terraform/live/accounts/dev/us-east-1/network
terragrunt init
terragrunt plan
terragrunt apply
```

📖 **[Complete Terraform Guide](../docs/infrastructure/TERRAFORM-GUIDE.md)**

## Quick Reference

| What | Command |
|------|---------|
| **Local Docker** | `make docker-up` |
| **Local SAM** | `cd infrastructure/sam && make local-dev` |
| **Deploy to AWS** | `cd infrastructure/sam && make deploy-all ENV=dev` |
| **Terraform Apply** | `cd terraform/live/accounts/<env>/<module> && terragrunt apply` |

## Documentation

- [Docker Guide](../docs/infrastructure/DOCKER.md) - Detailed Docker operations
- [SAM Guide](../docs/infrastructure/SAM-GUIDE.md) - SAM deployment and local testing
- [Terraform Guide](../docs/infrastructure/TERRAFORM-GUIDE.md) - Infrastructure as Code
- [SAM README](sam/README.md) - SAM-specific documentation
