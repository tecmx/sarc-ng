# Network Module Example

This example demonstrates how to create a complete VPC with public, private, and database subnets.

## Prerequisites

- AWS account with appropriate permissions
- Terraform >= 1.0.0
- AWS Provider >= 5.0.0

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- VPC with DNS support
- 2 public subnets (with NAT gateway)
- 2 private subnets
- 2 database subnets
- Internet Gateway
- NAT Gateway (single for dev)
- Route tables
- VPC Flow Logs (CloudWatch)

## Cleanup

```bash
terraform destroy
```

