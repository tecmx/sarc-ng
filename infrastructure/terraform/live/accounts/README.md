# Terragrunt Account-Based Structure

This directory contains the modernized Terragrunt account-based structure for SARC-NG infrastructure.

## Structure Overview

The structure follows modern Terragrunt best practices for organization:

```
accounts/
├── README.md                   # This file
├── account.hcl                 # Common account settings (if any)
├── dev/                        # Development account
│   ├── account.hcl             # Dev account-specific settings
│   ├── env.hcl                 # Dev environment variables
│   └── us-west-2/              # AWS Region
│       ├── region.hcl          # Region-specific settings
│       ├── network/            # Network resources
│       │   └── terragrunt.hcl
│       ├── compute/            # Compute resources
│       │   ├── ecs-cluster/
│       │   │   └── terragrunt.hcl
│       │   └── eks-cluster/
│       │       └── terragrunt.hcl
│       ├── database/           # Database resources
│       │   └── terragrunt.hcl
│       ├── dns/                # DNS resources
│       │   └── zone/
│       │       └── terragrunt.hcl
│       ├── observability/      # Monitoring resources
│       │   └── terragrunt.hcl
│       └── services/           # Application services
│           ├── sarcng-api/
│           │   ├── dns/record/
│           │   │   └── terragrunt.hcl
│           │   ├── ecs-alb-service/
│           │   │   └── terragrunt.hcl
│           │   └── schema/
│           │       └── terragrunt.hcl
│           └── sarcng-web/
│               ├── dns/record/
│               │   └── terragrunt.hcl
│               └── eks-helm-release/
│                   └── terragrunt.hcl
├── staging/                    # Staging account (similar structure)
└── prod/                       # Production account (similar structure)
```

## Key Features

1. **Account Isolation**: Each account (dev, staging, prod) is fully isolated with its own configuration.

2. **Configuration Sharing**: Common settings are shared through `account.hcl`, `region.hcl`, and `env.hcl` files.

3. **DRY Approach**: Using `include` with `expose = true` to share settings and local variables.

4. **Explicit Dependencies**: Clear dependency management between modules.

5. **Hooks**: Before/after hooks for validation, notification, and other tasks.

6. **Environment Variables**: Support for runtime environment variable overrides.

7. **LocalStack Support**: Built-in compatibility with LocalStack for local development.

## Usage

### Deploying an Entire Environment

To deploy all resources in an environment:

```bash
cd accounts/dev/us-west-2
terragrunt run-all init
terragrunt run-all plan
terragrunt run-all apply
```

### Deploying a Specific Module

To deploy a specific module:

```bash
cd accounts/dev/us-west-2/network
terragrunt init
terragrunt plan
terragrunt apply
```

### Using LocalStack for Local Development

For local development with LocalStack:

```bash
export LOCALSTACK=true
cd accounts/dev/us-west-2/network
terragrunt plan
```

## Best Practices

1. **Always use the include with expose**: Use `include { path = find_in_parent_folders() expose = true }` to share parent configurations.

2. **Manage dependencies explicitly**: Use `dependency` blocks with clear config_paths.

3. **Use consistent naming**: Follow the naming conventions for all resources.

4. **Keep environments isolated**: Each environment should be completely isolated from others.

5. **Add validation hooks**: Use hooks to validate inputs and configurations.

6. **Test in lower environments first**: Always test changes in dev before applying to staging or prod. 