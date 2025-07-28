# Infrastructure

Terraform and Terragrunt configuration for SARC-NG AWS infrastructure.

## Directory Structure

```
infrastructure/terraform/
├── README.md                  # This file
├── live/                      # Environment configurations
│   ├── terragrunt.hcl         # Root configuration
│   └── accounts/              # Account-based organization
│       ├── dev/               # Development environment
│       │   ├── account.hcl    # Account settings
│       │   ├── env.hcl        # Environment variables
│       │   ├── terragrunt.hcl # LocalStack configuration
│       │   └── us-east-1/     # Region deployments
│       │       ├── region.hcl # Region settings
│       │       ├── network/   # Network infrastructure
│       │       ├── database/  # Database infrastructure
│       │       └── compute/   # Compute infrastructure
│       ├── staging/           # Staging environment
│       └── prod/              # Production environment
└── modules/                   # Reusable Terraform modules
    ├── network/               # VPC, subnets, networking
    ├── compute/               # ECS, EKS, Lambda
    ├── database/              # RDS, Aurora
    ├── dns/                   # Route53
    └── observability/         # Monitoring, logging
```

## Quick Start

### Prerequisites

Required tools:
- Terraform >= 1.0.0
- Terragrunt >= 0.45.0
- AWS CLI
- LocalStack (for local development)

Installation:
```bash
# LocalStack
pip install localstack

# Core tools (macOS)
brew install terraform terragrunt awscli

# Core tools (Ubuntu)
sudo apt-get install terraform terragrunt awscli
```

### Basic Operations

Navigate to a module and run standard Terragrunt commands:

```bash
# Example: Network module in dev environment
cd live/accounts/dev/us-east-1/network

# Initialize
terragrunt init

# Plan changes
terragrunt plan

# Apply changes
terragrunt apply

# Show outputs
terragrunt output

# Destroy resources
terragrunt destroy
```

## Environment Operations

### Development with LocalStack

```bash
# Start LocalStack
localstack start -d

# Check status
curl -s http://localhost:4566/_localstack/health

# Apply to LocalStack
cd live/accounts/dev/us-east-1/network
LOCALSTACK=true terragrunt apply --auto-approve

# Destroy from LocalStack
LOCALSTACK=true terragrunt destroy --auto-approve
```

### AWS Environment Deployment

```bash
# Development environment
cd live/accounts/dev/us-east-1/network
terragrunt apply

# Staging environment
cd live/accounts/staging/us-east-1/network
terragrunt apply

# Production environment
cd live/accounts/prod/us-east-1/network
terragrunt apply
```

### Bulk Operations

Apply multiple modules at once:

```bash
# Apply all modules in an environment
cd live/accounts/dev/us-east-1
terragrunt run-all apply

# Plan all modules
terragrunt run-all plan

# Destroy all modules
terragrunt run-all destroy
```

## Module-Specific Commands

### Network Module

```bash
cd live/accounts/dev/us-east-1/network
terragrunt plan
terragrunt apply
```

### Database Module

```bash
cd live/accounts/dev/us-east-1/database
terragrunt plan
terragrunt apply
```

### Compute Module

```bash
cd live/accounts/dev/us-east-1/compute
terragrunt plan
terragrunt apply
```

## Workflows

### Local Development

1. Start LocalStack:
   ```bash
   localstack start -d
   ```

2. Test network module:
   ```bash
   cd live/accounts/dev/us-east-1/network
   LOCALSTACK=true terragrunt apply --auto-approve
   ```

3. Test other modules as needed

### Environment Promotion

1. **Development:**
   ```bash
   cd live/accounts/dev/us-east-1
   terragrunt run-all apply
   ```

2. **Staging:**
   ```bash
   cd live/accounts/staging/us-east-1
   terragrunt run-all plan
   terragrunt run-all apply
   ```

3. **Production:**
   ```bash
   cd live/accounts/prod/us-east-1
   terragrunt run-all plan
   terragrunt run-all apply
   ```

### Module Development

When creating new modules:

1. Create module under `modules/{module-name}/`
2. Add live configuration under `live/accounts/{account}/{region}/{module}/`
3. Test with LocalStack first
4. Deploy to dev, then staging, then production

## Troubleshooting

### Cache Issues

```bash
# Clean Terragrunt cache
find . -name ".terragrunt-cache" -type d -exec rm -rf {} + 2>/dev/null
find . -name ".terraform" -type d -exec rm -rf {} + 2>/dev/null
```

### LocalStack Issues

```bash
# Check LocalStack health
curl -s http://localhost:4566/_localstack/health

# Restart LocalStack
localstack start -d
```

### AWS Permission Issues

```bash
# Verify AWS credentials
aws sts get-caller-identity

# Configure AWS CLI
aws configure
```

### Debug Mode

```bash
# Verbose logging
cd live/accounts/dev/us-east-1/network
terragrunt plan --terragrunt-log-level debug
```

## Available Modules

### Network Module
- VPC with public and private subnets
- Internet Gateway and NAT Gateway
- Route tables and security groups
- Network ACLs

### Database Module
- RDS instances
- Parameter groups and subnet groups
- Security groups for database access
- Backup and monitoring configuration

### Compute Module
- ECS clusters and services
- Lambda functions
- Application Load Balancers
- Auto Scaling Groups

### DNS Module
- Route53 hosted zones
- DNS records
- SSL certificate management

### Observability Module
- CloudWatch log groups and metrics
- Alerting rules
- Monitoring dashboards

## Security

- Never commit AWS credentials
- Use IAM roles with least privilege
- Enable CloudTrail for audit logging
- Use separate AWS accounts per environment
- Store secrets in AWS Secrets Manager

## Configuration

### AWS Setup

```bash
# Configure AWS CLI
aws configure

# Set environment variables
export AWS_ACCESS_KEY_ID=your-key
export AWS_SECRET_ACCESS_KEY=your-secret
export AWS_DEFAULT_REGION=us-east-1
```

### LocalStack Setup

```bash
# For local development
export LOCALSTACK=true
```

## Utilities

### Code Formatting

```bash
# Format all Terraform files
terraform fmt -recursive
```

### Validation

```bash
# Validate configuration
cd live/accounts/dev/us-east-1/network
terragrunt validate
```

### Module Listing

```bash
# List available modules
ls live/accounts/dev/us-east-1/
``` 