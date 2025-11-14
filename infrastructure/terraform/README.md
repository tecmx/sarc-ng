# Terraform Operations Guide

Complete guide for Terraform infrastructure management.

## Directory Structure

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

## Basic Workflow

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

## Common Operations

### Format and Validate

```bash
# Format all Terraform files
terragrunt hclfmt

# Validate configuration
terragrunt validate
```

### State Management

```bash
# Show current state
terragrunt show

# List resources
terragrunt state list

# Refresh state
terragrunt refresh
```

### Targeted Operations

```bash
# Target specific resource
terragrunt apply -target=aws_vpc.main

# Auto-approve (use in CI/CD)
terragrunt apply -auto-approve

# Plan with output file
terragrunt plan -out=tfplan
terragrunt apply tfplan
```

## Troubleshooting

### Update Providers

```bash
terragrunt init -upgrade
```

### Force Unlock

```bash
# If state is stuck/locked
terraform force-unlock <lock-id>
```

### Common Issues

**State Lock**

If Terraform state is locked:

1. Wait for the lock to expire (usually 15 minutes)
2. If process died, force unlock with the lock ID shown in error
3. Never unlock if another process is actually running

**Module Changes**

After modifying modules:

```bash
terragrunt init -upgrade
```

**Provider Version Conflicts**

```bash
# Update to latest compatible versions
terragrunt init -upgrade

# Or lock to specific versions in versions.tf
```

