# EKS Namespace Module Example

Creates a Kubernetes namespace with resource quotas and network policies.

## Prerequisites

- Existing EKS cluster
- kubectl configured
- Kubernetes provider configured

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Kubernetes Namespace
- Resource Quotas
- Network Policy (deny-all by default)

## Cleanup

```bash
terraform destroy
```

