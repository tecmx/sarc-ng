/**
 * Observability module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

  # For QA environments with workspaces, add workspace suffix to name
  suffix = var.environment == "qa" && terraform.workspace != "default" ? "-${terraform.workspace}" : ""

  # Full name used for resource naming
  full_name = "${local.name}${local.suffix}"

  # Prometheus configuration
  prometheus_namespace = var.prometheus_namespace != "" ? var.prometheus_namespace : "monitoring"

  # Default Prometheus configuration
  default_prometheus_values = {
    server = {
      persistentVolume = {
        enabled      = true
        size         = var.prometheus_storage_size
        storageClass = var.storage_class_name
      }
      retention = "${var.prometheus_retention_period}d"
      resources = {
        limits = {
          cpu    = var.prometheus_cpu_limit
          memory = var.prometheus_memory_limit
        }
        requests = {
          cpu    = var.prometheus_cpu_request
          memory = var.prometheus_memory_request
        }
      }
    }
    alertmanager = {
      enabled = var.enable_alertmanager
      persistentVolume = {
        enabled      = var.enable_alertmanager
        size         = var.alertmanager_storage_size
        storageClass = var.storage_class_name
      }
    }
  }

  # Grafana configuration
  grafana_namespace = var.grafana_namespace != "" ? var.grafana_namespace : "monitoring"

  # Default Grafana configuration
  default_grafana_values = {
    persistence = {
      enabled      = true
      size         = var.grafana_storage_size
      storageClass = var.storage_class_name
    }
    adminPassword = var.grafana_admin_password
    resources = {
      limits = {
        cpu    = var.grafana_cpu_limit
        memory = var.grafana_memory_limit
      }
      requests = {
        cpu    = var.grafana_cpu_request
        memory = var.grafana_memory_request
      }
    }
    # Default datasource configuration for Prometheus
    datasources = {
      "datasources.yaml" = {
        apiVersion = 1
        datasources = [
          {
            name      = "Prometheus"
            type      = "prometheus"
            url       = "http://prometheus-server.${local.prometheus_namespace}.svc.cluster.local"
            access    = "proxy"
            isDefault = true
          }
        ]
      }
    }
    # Load standard dashboards
    dashboards = {
      default = true
    }
  }

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Component   = "Observability"
    },
    var.additional_tags
  )
} 
