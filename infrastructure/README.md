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

## Docker

### Commands

```bash
make docker-up       # Start services
make docker-down     # Stop services
make docker-logs     # View logs
make docker-clean    # Remove all data
```

### Services

- API: <http://localhost:8080/api/v1>
- Swagger: <http://localhost:8080/swagger/index.html>
- DB Admin: <http://localhost:8081>
- Metrics: <http://localhost:8080/metrics>

### Manual Commands

```bash
cd infrastructure/docker

# Start
docker compose up -d --wait

# Stop
docker compose down

# Stop and remove data
docker compose down -v

# Logs
docker compose logs -f app

# Database access
docker exec -it sarc-ng-db-dev mysql -u root -p
```

## AWS SAM (Lambda)

### Local Development (No AWS credentials needed)

```bash
cd infrastructure/sam

# Start everything (MySQL + API Gateway)
make local-dev

# Stop everything
make local-stop
```

**Configuration**: Edit `sam/local.env` for environment variables (simple key=value format).  
`env.json` is auto-generated from `local.env` automatically.

📖 **See [sam/README.md](sam/README.md) for detailed guide**

### Deployment to AWS (Requires AWS credentials)

```bash
cd infrastructure/sam

# Deploy
make deploy

# Show deployed URLs
make urls

# Delete stack
make delete
```

### Manual Commands

```bash
# Build
make build

# Clean build artifacts
make clean
```

---

## Terraform Commands

All commands run from the module directory.

### Directory Structure

```
infrastructure/terraform/live/accounts/
├── dev/
│   └── us-east-1/
│       ├── network/      # VPC, subnets, routing
│       ├── database/     # RDS, parameter groups
│       └── compute/      # ECS, Lambda, etc.
├── staging/
└── prod/
```

### Basic Workflow

```bash
# Navigate to module
cd infrastructure/terraform/live/accounts/dev/us-east-1/network

# Initialize (first time or after module changes)
terragrunt init

# Plan changes (preview)
terragrunt plan

# Apply changes
terragrunt apply

# Destroy resources (careful!)
terragrunt destroy
```

### Common Operations

```bash
# Format all Terraform files
terragrunt hclfmt

# Validate configuration
terragrunt validate

# Show current state
terragrunt show

# List resources
terragrunt state list

# Refresh state
terragrunt refresh

# Target specific resource
terragrunt apply -target=aws_vpc.main

# Auto-approve (use in CI/CD)
terragrunt apply -auto-approve

# Plan with output file
terragrunt plan -out=tfplan
terragrunt apply tfplan
  --docker-network sarc-ng-network \
  --env-vars env.json
```

## Troubleshooting

### Docker

```bash
# Check status
docker ps -a

# View logs
docker logs sarc-ng-server-dev

# Rebuild
docker compose down -v
docker compose build --no-cache
docker compose up -d --wait

# Database access
docker exec -it sarc-ng-db-dev mysql -u root -p
```

### SAM

```bash
# Validate template
sam validate

# Debug invoke
sam local invoke SarcNgFunction --event events/test.json --debug

# Check build
ls -la .aws-sam/build/SarcNgFunction/
```

### Terraform

```bash
# Update providers
terragrunt init -upgrade

# Force unlock (if stuck)
terraform force-unlock <lock-id>
```

## Quick Reference

### Local Development Workflows

| What | Command |
|------|---------|
| **Full local environment** | `cd infrastructure/sam && make local-dev` |
| **Just Docker services** | `cd infrastructure/docker && docker compose up -d` |
| **Stop local SAM** | `cd infrastructure/sam && make local-stop` |
| **Stop Docker services** | `cd infrastructure/docker && docker compose down` |

### AWS Deployment Workflows

| What | Command |
|------|---------|
| **Deploy to AWS** | `cd infrastructure/sam && make deploy` |
| **Show deployed URLs** | `cd infrastructure/sam && make urls` |
| **Delete AWS stack** | `cd infrastructure/sam && make delete` |
