# AWS SAM Deployment

AWS SAM deployment configuration for SARC-NG Lambda function.

## Quick Start

```bash
cd infrastructure/sam

# Validate template
make validate

# Deploy to AWS
make deploy
```

## Prerequisites

- AWS CLI configured with credentials
- AWS SAM CLI installed
- Go 1.21+ installed

## AWS Credentials

Edit `aws.env` with your AWS credentials:

```bash
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_SESSION_TOKEN=your-session-token
```

## Available Commands

```bash
make build        # Build Lambda function
make validate     # Validate SAM template
make deploy       # Deploy to AWS (requires aws.env)
make delete       # Delete CloudFormation stack
make status       # Show deployment status
make urls         # Show deployed URLs
make clean        # Remove build artifacts
make help         # Show all commands
```

## Deployment Workflow

1. **Configure AWS credentials** in `aws.env`
2. **Validate template**: `make validate`
3. **Deploy to AWS**: `make deploy`
4. **Check status**: `make status`
5. **Get URLs**: `make urls`

## Files

- `template.yaml` - CloudFormation/SAM template
- `samconfig.toml` - SAM deployment configuration
- `aws.env` - AWS credentials 🔐 **Edit this**
- `Makefile` - Build and deployment commands

## Deployed Resources

After deployment, the stack creates:

- **Lambda Function** - SARC-NG API handler
- **API Gateway** - REST API endpoint
- **VPC** - Isolated network for Lambda and RDS
- **RDS MySQL** - Database instance
- **Secrets Manager** - Database credentials storage

Access your deployed app via the URLs from `make urls`.

## Troubleshooting

### AWS credentials not set

Ensure `aws.env` contains valid AWS credentials with proper permissions.

### Template validation fails

Run `make validate` to check for syntax errors in `template.yaml`.

### Deployment fails

Check CloudFormation console for detailed error messages or run:

```bash
aws cloudformation describe-stack-events --stack-name sarc-ng-prod
```
