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

- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html
- DB Admin: http://localhost:8081
- Metrics: http://localhost:8080/metrics

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

```bash
cd infrastructure/sam

# Build
sam build

# Deploy (first time)
sam deploy --guided

# Deploy (subsequent)
sam build && sam deploy

# Local testing
sam local start-api --port 3001 \
  --docker-network sarc-ng-network \
  --env-vars env.json
```

### Local Testing Setup

```bash
# Start database
cd ../docker && docker compose up -d --wait

# Create env file
cd ../sam
cp env.json.example env.json
# Edit env.json with database credentials

# Start local API
sam build
sam local start-api --port 3001 \
sam deploy \
  --parameter-overrides \
    Environment=dev \
    DBPassword=secure-password
```

### Validate Template

```bash
sam validate
```

Checks `template.yaml` syntax and configuration.

### Delete Stack

```bash
sam delete
```

Deletes the entire CloudFormation stack and all AWS resources.

### Get Stack Outputs

```bash
# Using SAM
sam list stack-outputs

# Using AWS CLI
aws cloudformation describe-stacks \
  --stack-name sarc-ng-dev \
  --query 'Stacks[0].Outputs'
```

### Clean Build Artifacts

```bash
# If writable
rm -rf .aws-sam

# If Docker-created (requires sudo)
sudo rm -rf .aws-sam
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

### Cleanup

```bash
# Delete SAM stack
sam delete --stack-name <stack-name>
```

## Terraform

```bash
cd infrastructure/terraform/live/accounts/dev/us-east-1/<module>

# Plan
terragrunt plan

# Apply
terragrunt apply

# Multi-environment
cd infrastructure/terraform/live/accounts/prod/us-east-1/<module>
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

## Common Workflows

### Local Development
```bash
cd infrastructure/docker
docker compose up -d --wait
# Access: http://localhost:8080
docker compose down
```

### Lambda Testing
```bash
cd infrastructure/docker && docker compose up -d --wait
cd ../sam && sam build
sam local start-api --port 3001 --docker-network sarc-ng-network --env-vars env.json
```

### AWS Deployment
```bash
cd infrastructure/sam
sam build && sam validate && sam deploy
```
