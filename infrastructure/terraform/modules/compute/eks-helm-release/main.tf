/**
 * EKS Helm Release module - Deploys applications via Helm charts on EKS
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "cluster_endpoint" {
  description = "Endpoint of the EKS cluster"
  type        = string
}

variable "cluster_ca_data" {
  description = "CA certificate data for the EKS cluster"
  type        = string
}

variable "service_name" {
  description = "Name of the service"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace to deploy the chart to"
  type        = string
  default     = null
}

variable "create_namespace" {
  description = "Whether to create the namespace"
  type        = bool
  default     = true
}

variable "release_name" {
  description = "Name of the Helm release"
  type        = string
  default     = null
}

variable "chart_repo" {
  description = "Repository URL where to locate the Helm chart"
  type        = string
  default     = null
}

variable "chart" {
  description = "Chart name to be installed"
  type        = string
}

variable "chart_version" {
  description = "Version of the Helm chart"
  type        = string
  default     = null
}

variable "values" {
  description = "List of values in raw YAML to pass to helm"
  type        = list(string)
  default     = []
}

variable "set" {
  description = "Value block with custom values to be merged with the values yaml"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "set_sensitive" {
  description = "Value block with sensitive values to be merged with the values yaml"
  type = list(object({
    name  = string
    value = string
  }))
  default   = []
  sensitive = true
}

variable "create_service_account" {
  description = "Whether to create a service account for the Helm release"
  type        = bool
  default     = false
}

variable "service_account_name" {
  description = "Name of the service account to create"
  type        = string
  default     = null
}

variable "service_account_annotations" {
  description = "Annotations to add to the service account"
  type        = map(string)
  default     = {}
}

variable "timeout" {
  description = "Time in seconds to wait for any individual Kubernetes operation"
  type        = number
  default     = 300
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
}

locals {
  name = "${var.project_name}-${var.environment}"

  # For QA environments with workspaces, add workspace suffix to names
  service_name_suffix = var.environment == "qa" && terraform.workspace != "default" ? "-${terraform.workspace}" : ""
  service_name        = "${var.service_name}${local.service_name_suffix}"

  # Default release name and namespace
  release_name = var.release_name != null ? var.release_name : local.service_name
  namespace    = var.namespace != null ? var.namespace : local.service_name

  # Service account name
  service_account_name = var.service_account_name != null ? var.service_account_name : local.service_name

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Service     = local.service_name
    },
    var.additional_tags
  )
}

# Configure the Kubernetes provider
provider "kubernetes" {
  host                   = var.cluster_endpoint
  cluster_ca_certificate = base64decode(var.cluster_ca_data)

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "aws"
    args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
  }
}

# Configure the Helm provider
provider "helm" {
  kubernetes {
    host                   = var.cluster_endpoint
    cluster_ca_certificate = base64decode(var.cluster_ca_data)

    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
    }
  }
}

# Create namespace if required
resource "kubernetes_namespace" "this" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name = local.namespace

    labels = {
      name        = local.namespace
      project     = var.project_name
      environment = var.environment
      managed-by  = "Terraform"
    }
  }
}

# Create service account if required
resource "kubernetes_service_account" "this" {
  count = var.create_service_account ? 1 : 0

  metadata {
    name        = local.service_account_name
    namespace   = var.create_namespace ? kubernetes_namespace.this[0].metadata[0].name : local.namespace
    annotations = var.service_account_annotations

    labels = {
      app         = local.service_name
      project     = var.project_name
      environment = var.environment
      managed-by  = "Terraform"
    }
  }

  automount_service_account_token = true

  depends_on = [kubernetes_namespace.this]
}

# Deploy Helm chart
resource "helm_release" "this" {
  name             = var.name
  namespace        = var.namespace
  create_namespace = var.create_namespace
  repository       = var.repository
  chart            = var.chart
  version          = var.chart_version
  values           = var.values
  timeout          = var.timeout

  dynamic "set" {
    for_each = var.set
    content {
      name  = set.value.name
      value = set.value.value
    }
  }

  dynamic "set_sensitive" {
    for_each = var.set_sensitive
    content {
      name  = set_sensitive.value.name
      value = set_sensitive.value.value
    }
  }

  dynamic "postrender" {
    for_each = length(var.postrender) > 0 ? [var.postrender[0]] : []
    content {
      binary_path = postrender.value.binary_path
      args        = postrender.value.args
    }
  }

  atomic                = var.atomic
  cleanup_on_fail       = var.cleanup_on_fail
  dependency_update     = var.dependency_update
  disable_webhooks      = var.disable_webhooks
  force_update          = var.force_update
  max_history           = var.max_history
  recreate_pods         = var.recreate_pods
  render_subchart_notes = var.render_subchart_notes
  replace               = var.replace
  reuse_values          = var.reuse_values
  skip_crds             = var.skip_crds
  wait                  = var.wait
  wait_for_jobs         = var.wait_for_jobs
  description           = var.description
}

output "namespace" {
  description = "The namespace where the release was deployed"
  value       = var.create_namespace ? kubernetes_namespace.this[0].metadata[0].name : local.namespace
}

output "release_name" {
  description = "The name of the Helm release"
  value       = helm_release.this.name
}

output "chart_version" {
  description = "The version of the chart that was deployed"
  value       = helm_release.this.version
}

output "service_account_name" {
  description = "The name of the service account"
  value       = var.create_service_account ? kubernetes_service_account.this[0].metadata[0].name : null
}
