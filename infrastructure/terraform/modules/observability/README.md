# Observability Module

This module deploys a complete observability stack for Kubernetes clusters, including:

- Prometheus for metrics collection and alerting
- Grafana for visualization and dashboards

## Features

- Prometheus server with configurable retention and storage
- AlertManager for alert handling (optional)
- Grafana with pre-configured dashboards
- Ingress support for Grafana
- Default alerting rules for common scenarios
- Customizable resource settings for different environments

## Usage

```hcl
module "observability" {
  source = "../../modules/observability"
  
  # EKS cluster configuration
  cluster_endpoint       = module.eks_cluster.cluster_endpoint
  cluster_ca_certificate = module.eks_cluster.cluster_certificate_authority_data
  cluster_name           = module.eks_cluster.cluster_name

  # Environment configuration
  project_name = "sarc-ng"
  environment  = "dev"
  
  # Prometheus configuration
  prometheus_retention_period = 15
  prometheus_storage_size     = "10Gi"
  
  # Grafana configuration
  grafana_storage_size    = "5Gi"
  grafana_admin_password  = "secure-password"
  
  # Ingress configuration
  enable_ingress          = true
  ingress_host            = "grafana.example.com"
  ingress_tls_secret_name = "grafana-tls"
  ingress_annotations     = {
    "kubernetes.io/ingress.class" = "nginx"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| project_name | Project name for resource naming | `string` | n/a | yes |
| environment | Environment name (dev, qa, staging, prod) | `string` | n/a | yes |
| cluster_endpoint | EKS cluster endpoint | `string` | n/a | yes |
| cluster_ca_certificate | EKS cluster CA certificate | `string` | n/a | yes |
| cluster_token | EKS cluster authentication token | `string` | `null` | no |
| cluster_name | EKS cluster name | `string` | n/a | yes |
| storage_class_name | Kubernetes storage class name | `string` | `"gp2"` | no |
| prometheus_namespace | Kubernetes namespace for Prometheus | `string` | `"monitoring"` | no |
| prometheus_helm_chart_version | Prometheus Helm chart version | `string` | `"15.10.1"` | no |
| prometheus_storage_size | Storage size for Prometheus server | `string` | `"8Gi"` | no |
| prometheus_retention_period | Data retention period in days | `number` | `15` | no |
| enable_alertmanager | Whether to enable AlertManager | `bool` | `true` | no |
| grafana_namespace | Kubernetes namespace for Grafana | `string` | `"monitoring"` | no |
| grafana_helm_chart_version | Grafana Helm chart version | `string` | `"6.50.0"` | no |
| grafana_storage_size | Storage size for Grafana | `string` | `"10Gi"` | no |
| grafana_admin_password | Admin password for Grafana | `string` | `"admin"` | no |
| enable_ingress | Whether to create Ingress for Grafana | `bool` | `false` | no |
| ingress_host | Hostname for Grafana Ingress | `string` | `""` | no |
| ingress_annotations | Annotations for Grafana Ingress | `map(string)` | `{}` | no |
| ingress_tls_secret_name | TLS secret for Grafana Ingress | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| prometheus_namespace | Namespace where Prometheus is deployed |
| grafana_namespace | Namespace where Grafana is deployed |
| prometheus_server_endpoint | Internal endpoint for Prometheus server |
| alertmanager_endpoint | Internal endpoint for AlertManager |
| grafana_endpoint | Internal endpoint for Grafana |
| grafana_ingress_host | Hostname for Grafana ingress |

## Accessing Grafana

After deploying the module:

1. If ingress is enabled, access Grafana via the configured hostname
2. If ingress is not enabled, use port forwarding:
   ```
   kubectl port-forward svc/grafana 3000:80 -n monitoring
   ```
3. Login with username `admin` and the password specified in `grafana_admin_password`

## Customizing Alerts

To customize alerts, provide your own alert rules:

```hcl
module "observability" {
  # ... other configuration

  alert_rules = {
    "custom-alerts.yaml" = yamlencode({
      groups = [
        {
          name = "custom-alerts"
          rules = [
            {
              alert = "HighErrorRate"
              expr  = "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m])) * 100 > 10"
              for   = "5m"
              labels = {
                severity = "critical"
              }
              annotations = {
                summary     = "High error rate detected"
                description = "Error rate is above 10% (current value: {{ $value }}%)"
              }
            }
          ]
        }
      ]
    })
  }
} 