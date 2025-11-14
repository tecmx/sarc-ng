# EKS Helm Release Module Example

Deploys a Helm chart to an EKS cluster.

## Prerequisites

- Existing EKS cluster
- Helm 3.x installed
- kubectl configured

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Helm Release
- Kubernetes resources defined by the chart

## Cleanup

```bash
terraform destroy
```

