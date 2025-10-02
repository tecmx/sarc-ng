/**
 * EKS Namespace module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

  # Service account name - use specified name or namespace name
  service_account_name = var.service_account_name != "" ? var.service_account_name : var.namespace

  # For QA environments with workspaces, add workspace suffix to namespace
  namespace_suffix = var.environment == "qa" && terraform.workspace != "default" ? "-${terraform.workspace}" : ""
  namespace_name   = "${var.namespace}${local.namespace_suffix}"

  # Default labels for namespace
  default_labels = {
    "app.kubernetes.io/managed-by" = "terraform"
    "app.kubernetes.io/part-of"    = var.project_name
    "environment"                  = var.environment
  }

  # Merge default and custom labels
  labels = merge(local.default_labels, var.labels)

  # Default annotations for namespace
  default_annotations = {}

  # Merge default and custom annotations
  annotations = merge(local.default_annotations, var.annotations)

  # Tags for AWS resources
  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Namespace   = local.namespace_name
    },
    var.additional_tags
  )
} 
