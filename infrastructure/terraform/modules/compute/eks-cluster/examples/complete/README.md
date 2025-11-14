# EKS Cluster Module Example

Creates an EKS cluster with managed node groups and IRSA enabled.

## Prerequisites

- Existing VPC with private subnets
- Terraform >= 1.0.0
- kubectl installed

## Usage

```bash
terraform init
terraform plan
terraform apply

# Configure kubectl
aws eks update-kubeconfig --region us-east-1 --name sarc-ng-dev-eks
```

## Resources Created

- EKS Cluster
- Managed node groups
- IAM roles and policies
- Security groups
- OIDC provider (for IRSA)
- SSM parameters

## Cleanup

```bash
terraform destroy
```

