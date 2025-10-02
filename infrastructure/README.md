# Infrastructure

SARC-NG infrastructure management - Docker, AWS SAM, and Terraform configurations.

## Quick Start

```bash
# Start local development environment
make docker-up        # From project root
# or
cd infrastructure && make docker-up
```

## Prerequisites

- **Docker** - For local development
- **AWS CLI** - For SAM and Terraform deployments (configured with credentials)
- **AWS SAM CLI** - For Lambda deployment and testing (install: `pip install aws-sam-cli` or `brew install aws-sam-cli`)
- **Terraform/Terragrunt** - For infrastructure provisioning (optional)

---

## Docker Commands

### Using Makefile (Recommended)

From project root or `infrastructure/` directory:

```bash
make docker-up       # Start services and wait for health
make docker-down     # Stop services (keeps data)
make docker-logs     # View all logs
make docker-logs service=app    # View specific service logs
make docker-clean    # Remove all data (WARNING: deletes database)
```

### Using Docker Compose Directly

All commands run from `infrastructure/docker` directory:

```bash
cd infrastructure/docker
```

#### Start Services

```bash
# Start all services and wait for health checks
docker compose up -d --wait

# Start with build (after Dockerfile changes)
docker compose up -d --build --wait
```

**Services started:**

- API: <http://localhost:8080/api/v1>
- Swagger: <http://localhost:8080/swagger/index.html>
- DB Admin: <http://localhost:8081>
- Metrics: <http://localhost:8080/metrics>
- Debug Port: localhost:2345

### Stop Services

```bash
# Stop services (keeps data)
docker compose down

# Stop and remove all data (WARNING: deletes database!)
docker compose down -v
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f db
docker compose logs -f app

# Last 100 lines
docker compose logs --tail 100
```

### Service Management

```bash
# Check status
docker compose ps

# Restart specific service
docker compose restart app

# Rebuild and restart
docker compose up -d --build app
```

### Database Access

```bash
# MySQL CLI
docker exec -it sarc-ng-db-dev mysql -u root -p
# Password: example

# Or use Adminer web UI
open http://localhost:8081
```

---

## SAM (Lambda) Commands

All commands run from `infrastructure/sam` directory:

```bash
cd infrastructure/sam
```

### Build Lambda

```bash
sam build
```

Compiles Go binary for Lambda deployment:

- Target: `GOOS=linux GOARCH=amd64`
- Output: `.aws-sam/build/SarcNgFunction/bootstrap`

### Local Testing

**Prerequisites:** Database must be running

```bash
# Start database if not running
cd ../docker && docker compose up -d --wait && cd ../sam

# Create environment file (first time only)
cp env.json.example env.json
# Edit env.json with your database credentials

# Start local API (requires build first)
sam local start-api \
  --port 3001 \
  --docker-network sarc-ng-network \
  --env-vars env.json
```

Local API will be available at: <http://localhost:3001>

**Test specific function:**

```bash
sam local invoke SarcNgFunction --event events/test-event.json
```

**How it works:**

- Lambda runs in Docker container
- Connects to `sarc-ng-network` to access database at `db:3306`
- Environment variables loaded from `env.json`
- Press `Ctrl+C` to stop
- Hot reload: code changes require `sam build` to rebuild

### Deploy to AWS

**First time deployment:**

```bash
sam deploy --guided
```

Answer prompts:

- Stack Name: `sarc-ng-dev`
- AWS Region: `us-east-1`
- Confirm changes: `Y`
- Allow SAM CLI IAM role creation: `Y`
- Save arguments to config: `Y`

**Subsequent deployments:**

```bash
sam build
sam deploy
```

Uses saved configuration from `samconfig.toml`.

**Environment-specific deployment:**

```bash
sam deploy --config-env development
sam deploy --config-env production
```

**Custom parameters:**

```bash
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
```

### Multi-Environment

```bash
# Development
cd infrastructure/terraform/live/accounts/dev/us-east-1/network
terragrunt plan

# Staging
cd infrastructure/terraform/live/accounts/staging/us-east-1/network
terragrunt plan

# Production
cd infrastructure/terraform/live/accounts/prod/us-east-1/network
terragrunt plan
```

---

## Configuration Files

### Docker Configuration

**File:** `docker/docker-compose.yml`

**Services:**

- `app` - Go application server (port 8080)
- `db` - MySQL 8.0 database (port 3306)
- `adminer` - Database admin UI (port 8081)

**Environment Variables:**

- `DB_HOST` - Database host (default: db)
- `DB_PORT` - Database port (default: 3306)
- `DB_USER` - Database user (default: root)
- `DB_PASSWORD` - Database password (default: example)
- `DB_NAME` - Database name (default: sarcng)

**Volumes:**

- `db_data` - Persistent database storage
- `go-modules` - Go module cache
- `go-cache` - Go build cache

### SAM Configuration

**File:** `sam/env.json` (gitignored)

Local environment overrides for `sam local`:

```json
{
  "SarcNgFunction": {
    "DB_HOST": "db",
    "DB_PORT": "3306",
    "DB_USER": "root",
    "DB_PASSWORD": "example",
    "DB_NAME": "sarcng",
    "GIN_MODE": "debug",
    "LOG_LEVEL": "debug",
    "ENVIRONMENT": "dev"
  }
}
```

**File:** `sam/template.yaml`

**Resources:**

- `SarcNgFunction` - Lambda function (Go provided.al2 runtime)
- `SarcDatabase` - RDS MySQL instance
- `ApiGateway` - HTTP API endpoint

**Template Parameters:**

- `Environment` - Deployment environment (dev/staging/prod)
- `DBPassword` - RDS MySQL password
- `DBName` - Database name (default: sarcng)

**Lambda Environment Variables:**

- `DB_HOST` - RDS endpoint (auto-configured)
- `DB_PORT` - Database port (3306)
- `DB_NAME` - Database name
- `DB_USER` - Database user (root)
- `DB_PASSWORD` - From parameter
- `GIN_MODE` - release
- `LOG_LEVEL` - info
- `ENVIRONMENT` - From parameter

**File:** `sam/samconfig.toml`

Generated after first deployment. Contains:

- Stack name
- AWS region
- S3 bucket for artifacts
- CloudFormation capabilities
- Parameter overrides

### Terraform Configuration

**Structure:**

- `modules/` - Reusable Terraform modules
- `live/accounts/<env>/<region>/<module>/` - Environment-specific configurations
- `terragrunt.hcl` - Terragrunt configuration

**Backend:** S3 with DynamoDB state locking (configured per environment)

---

## Debugging

### Docker Issues

```bash
# Check container status
docker ps -a

# View container logs
docker logs sarc-ng-db-dev
docker logs sarc-ng-server-dev

# Inspect container
docker inspect sarc-ng-db-dev

# Check networks
docker network ls
docker network inspect sarc-ng-network

# Check volumes
docker volume ls
docker volume inspect docker_db_data

# Execute command in container
docker exec -it sarc-ng-server-dev sh

# Rebuild from scratch
cd infrastructure/docker
docker compose down -v
docker compose build --no-cache
docker compose up -d --wait
```

### SAM Issues

```bash
# Check build output
ls -la .aws-sam/build/SarcNgFunction/

# Verify binary
file .aws-sam/build/SarcNgFunction/bootstrap

# Check Lambda container logs
docker logs <container-name>

# Find Lambda container
docker ps --filter "ancestor=public.ecr.aws/lambda/provided:al2-rapid-x86_64"

# Test with verbose logging
sam local invoke SarcNgFunction \
  --event events/test.json \
  --debug

# Validate template
sam validate --debug
```

### Terraform Issues

**State lock issues:**

```bash
# Check for stuck locks
# Manually release if needed (with caution)
terraform force-unlock <lock-id>
```

**Provider version conflicts:**

```bash
# Update providers
terragrunt init -upgrade
```

**Module not found:**

```bash
# Ensure in correct directory
cd infrastructure/terraform/live/accounts/dev/us-east-1/<module>

# Re-initialize
terragrunt init
```

---

## Lambda Local Development with Database

**Industry Standard Approach**: When running Lambda functions locally, the Lambda runs inside a Docker container. To access services on your host (like databases), SAM CLI provides the `--docker-network` flag to connect the Lambda container to the same Docker network as your other services.

**How it works:**

1. Docker Compose creates `sarc-ng-network`
2. Database container joins this network with hostname `db`
3. SAM local starts Lambda container
4. Lambda container also joins `sarc-ng-network` via `--docker-network` flag
5. Lambda can reach database at `db:3306` using Docker DNS

**Why this approach:**

| Approach | Works | Portable | Industry Standard |
|----------|-------|----------|-------------------|
| `host.docker.internal` | ❌ Linux | ❌ | No |
| Bridge IP (172.17.0.1) | ⚠️ Fragile | ❌ | No |
| **Docker Network** | ✅ | ✅ | **✅ AWS Recommended** |

---

## Best Practices

1. **Use Direct Commands** - Transparent and easy to customize
2. **Validate Before Deploy** - Always run `sam validate` or `terragrunt plan` first
3. **Clean Regularly** - Remove `.aws-sam` and old Terraform state
4. **Monitor Costs** - Check AWS billing after deployments
5. **Version Control** - Commit `samconfig.toml` and Terraform configs (not state)
6. **Security** - Never commit credentials or `env.json` files
7. **Use `--wait`** - Always use `docker compose up -d --wait` to ensure services are healthy

---

## Common Workflows

### Local Development

```bash
# 1. Start database and services
cd infrastructure/docker
docker compose up -d --wait

# 2. Access API
curl http://localhost:8080/health

# 3. View logs
docker compose logs -f app

# 4. Stop when done
docker compose down
```

### Lambda Development

```bash
# 1. Start database
cd infrastructure/docker
docker compose up -d --wait

# 2. Build Lambda
cd ../sam
sam build

# 3. Test locally
sam local start-api \
  --port 3001 \
  --docker-network sarc-ng-network \
  --env-vars env.json

# 4. Test endpoint
curl http://localhost:3001/health
```

### Deploy to AWS

```bash
# 1. Build Lambda
cd infrastructure/sam
sam build

# 2. Validate
sam validate

# 3. Deploy
sam deploy

# 4. Test deployed API
curl https://your-api-id.execute-api.us-east-1.amazonaws.com/Prod/health
```

### Infrastructure Changes

```bash
# 1. Navigate to module
cd infrastructure/terraform/live/accounts/dev/us-east-1/network

# 2. Make changes to *.tf files

# 3. Plan changes
terragrunt plan

# 4. Review output carefully

# 5. Apply if looks good
terragrunt apply
```

---

## Further Reading

- **AWS SAM:** <https://docs.aws.amazon.com/serverless-application-model/>
- **Docker Compose:** <https://docs.docker.com/compose/>
- **Terraform:** <https://www.terraform.io/docs>
- **Terragrunt:** <https://terragrunt.gruntwork.io/>
