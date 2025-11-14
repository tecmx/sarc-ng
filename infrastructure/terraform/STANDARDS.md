# Terraform Infrastructure Standards

All modules have been standardized to follow consistent patterns, security best practices, and modern Terraform conventions.

## Version Requirements

- **Terraform**: `>= 1.0.0` (all modules)
- **AWS Provider**: `>= 5.0.0` (all modules)
- **External Modules**: Pinned with `~>` syntax

## Code Standards

### File Organization

All modules use consistent file structure:
```
module/
├── versions.tf      # Terraform and provider requirements
├── variables.tf     # Input variables with validation
├── locals.tf        # Local values and computed tags
├── main.tf          # Primary resources
├── outputs.tf       # Output values
├── *.tf             # Additional logical groupings (iam.tf, security_groups.tf, etc.)
└── examples/complete/  # Working example usage
    ├── main.tf
    └── README.md
```

### Variable Naming

**Standard variables (all modules):**
- `project_name` - Project identifier (validated: lowercase-hyphen only)
- `environment` - Environment name (validated: dev, qa, staging, prod)
- `additional_tags` - Optional custom tags (map)

### Tagging

All resources automatically tagged with:
```hcl
{
  Project     = var.project_name
  Environment = var.environment
  ManagedBy   = "Terraform"
}
```

Custom tags merged via `additional_tags` variable.

### Input Validation

All modules validate:
- `project_name`: `^[a-z0-9-]+$` (lowercase alphanumeric with hyphens)
- `environment`: Must be one of: `dev`, `qa`, `staging`, `prod`

## Security Standards

### Security Groups

- **Pattern**: Separate `aws_security_group_rule` resources (not inline)
- **Lifecycle**: `create_before_destroy = true` on all security groups
- **Default**: No permissive defaults (0.0.0.0/0)

### Database Module

Required variables:
- `vpc_cidr` - Egress restricted to VPC CIDR only
- `allowed_cidr_blocks` or `allowed_security_group_ids` - Explicit ingress rules

### ECS NLB Service Module

Required variables:
- `allowed_cidr_blocks` - Must specify allowed CIDR blocks (validated: non-empty list)

### Lambda Modules

Security features:
- Optional X-Ray tracing
- Optional VPC configuration
- CloudWatch Logs with configurable retention
- IAM roles with least-privilege policies

## Module-Specific Standards

### Lambda Modules

Split into logical files:
- `main.tf` - Lambda function resources
- `iam.tf` - IAM roles and policies
- `api_gateway.tf` - API Gateway resources (http-api module)

### ECS Modules

- Fargate capacity providers with cost optimization defaults
- Container Insights optional
- Auto-scaling capabilities

### EKS Modules

- IRSA enabled by default
- Managed node groups
- Standard add-ons via separate modules

### Network Module

- Multi-AZ by default
- Separate subnet tiers (public, private, database)
- Flow logs optional
- Kubernetes tagging ready

## Examples

Every module includes `examples/complete/` with:
- Working configuration
- README with prerequisites and usage
- Realistic parameters
- Cleanup instructions

## Terragrunt Structure

```
live/
├── terragrunt.hcl              # Root config (remote state, provider, common inputs)
└── accounts/
    └── {env}/
        ├── account.hcl          # Account-specific config
        ├── env.hcl              # Environment variables
        ├── terragrunt.hcl       # Environment overrides (if needed)
        └── {region}/
            ├── region.hcl       # Region config
            └── {module}/
                └── terragrunt.hcl  # Module-specific config
```

## Usage Pattern

Standard module usage:
```hcl
terraform {
  source = "${get_path_to_repo_root()}/infrastructure/terraform/modules/{module-name}"
}

dependency "network" {
  config_path = "../network"
}

inputs = {
  # Standard variables (required)
  project_name = "sarc-ng"
  environment  = "dev"

  # Module-specific variables
  vpc_id = dependency.network.outputs.vpc_id
  # ...

  # Optional custom tags
  additional_tags = {
    Team = "Engineering"
  }
}
```

## Validation

Use standard Terraform tools:
```bash
# Format
terragrunt hclfmt

# Validate
terragrunt validate

# Plan
terragrunt plan

# Apply
terragrunt apply
```

## Reference

- Module examples: `modules/*/examples/complete/`
- Operations guide: `README.md`
- Each module's `variables.tf` for complete variable documentation

