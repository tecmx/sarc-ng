# Observability Module Example

Deploys Prometheus, Grafana, and monitoring stack to EKS cluster.

## Prerequisites

- Existing EKS cluster
- kubectl configured
- Sufficient cluster resources

## Usage

```bash
terraform init
terraform plan
terraform apply

# Access Grafana
kubectl port-forward -n monitoring svc/grafana 3000:80
# Navigate to http://localhost:3000
```

## Resources Created

- Prometheus StatefulSet
- Grafana Deployment
- Alert Manager
- CloudWatch dashboards
- Service monitors

## Cleanup

```bash
terraform destroy
```

