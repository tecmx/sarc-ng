# ECS Cluster Module Example

Creates an ECS cluster with Fargate capacity providers and Container Insights.

## Prerequisites

- AWS account with ECS permissions
- Terraform >= 1.0.0

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- ECS Cluster
- Fargate and Fargate Spot capacity providers
- CloudWatch Container Insights
- SSM parameters for cluster configuration

## Cleanup

```bash
terraform destroy
```

