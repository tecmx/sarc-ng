/**
 * Observability module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

# Kubernetes configuration
variable "cluster_endpoint" {
  description = "EKS cluster endpoint"
  type        = string
}

variable "cluster_ca_certificate" {
  description = "EKS cluster CA certificate"
  type        = string
}

variable "cluster_token" {
  description = "EKS cluster authentication token"
  type        = string
  sensitive   = true
  default     = null
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "storage_class_name" {
  description = "Kubernetes storage class name for persistent volumes"
  type        = string
  default     = "gp2"
}

# Prometheus configuration
variable "prometheus_namespace" {
  description = "Kubernetes namespace for Prometheus"
  type        = string
  default     = "monitoring"
}

variable "prometheus_helm_chart_version" {
  description = "Prometheus Helm chart version"
  type        = string
  default     = "15.10.1" # Set to a specific version for stability
}

variable "prometheus_helm_chart_repository" {
  description = "Prometheus Helm chart repository"
  type        = string
  default     = "https://prometheus-community.github.io/helm-charts"
}

variable "prometheus_storage_size" {
  description = "Storage size for Prometheus server"
  type        = string
  default     = "8Gi"
}

variable "prometheus_retention_period" {
  description = "Data retention period in days"
  type        = number
  default     = 15
}

variable "prometheus_cpu_request" {
  description = "CPU request for Prometheus"
  type        = string
  default     = "200m"
}

variable "prometheus_cpu_limit" {
  description = "CPU limit for Prometheus"
  type        = string
  default     = "1000m"
}

variable "prometheus_memory_request" {
  description = "Memory request for Prometheus"
  type        = string
  default     = "512Mi"
}

variable "prometheus_memory_limit" {
  description = "Memory limit for Prometheus"
  type        = string
  default     = "2Gi"
}

variable "prometheus_scrape_configs" {
  description = "Additional scrape configurations for Prometheus"
  type        = any
  default     = {}
}

variable "enable_alertmanager" {
  description = "Whether to enable the Prometheus Alertmanager"
  type        = bool
  default     = true
}

variable "alertmanager_storage_size" {
  description = "Storage size for Alertmanager"
  type        = string
  default     = "2Gi"
}

variable "alert_rules" {
  description = "Map of Prometheus alerting rules"
  type        = any
  default     = {}
}

# Grafana configuration
variable "grafana_namespace" {
  description = "Kubernetes namespace for Grafana"
  type        = string
  default     = "monitoring"
}

variable "grafana_helm_chart_version" {
  description = "Grafana Helm chart version"
  type        = string
  default     = "6.50.0" # Set to a specific version for stability
}

variable "grafana_helm_chart_repository" {
  description = "Grafana Helm chart repository"
  type        = string
  default     = "https://grafana.github.io/helm-charts"
}

variable "grafana_storage_size" {
  description = "Storage size for Grafana"
  type        = string
  default     = "10Gi"
}

variable "grafana_admin_password" {
  description = "Admin password for Grafana"
  type        = string
  sensitive   = true
  default     = "admin" # Should be overridden in production
}

variable "grafana_cpu_request" {
  description = "CPU request for Grafana"
  type        = string
  default     = "100m"
}

variable "grafana_cpu_limit" {
  description = "CPU limit for Grafana"
  type        = string
  default     = "500m"
}

variable "grafana_memory_request" {
  description = "Memory request for Grafana"
  type        = string
  default     = "128Mi"
}

variable "grafana_memory_limit" {
  description = "Memory limit for Grafana"
  type        = string
  default     = "512Mi"
}

variable "grafana_additional_data_sources" {
  description = "Additional data sources for Grafana"
  type        = any
  default     = {}
}

variable "grafana_additional_plugins" {
  description = "Additional plugins to install in Grafana"
  type        = list(string)
  default     = ["grafana-piechart-panel", "grafana-worldmap-panel"]
}

# Service exposure configuration
variable "enable_ingress" {
  description = "Whether to create an Ingress for Grafana"
  type        = bool
  default     = false
}

variable "ingress_host" {
  description = "Hostname for the Grafana Ingress"
  type        = string
  default     = ""
}

variable "ingress_annotations" {
  description = "Annotations for the Grafana Ingress"
  type        = map(string)
  default     = {}
}

variable "ingress_tls_secret_name" {
  description = "Name of the TLS secret for the Grafana Ingress"
  type        = string
  default     = ""
}

# Common tags
variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
