/**
 * ECS ALB Service module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}-${var.service_name}"

  # For QA environments with workspaces, add workspace suffix to name
  service_name_suffix = var.environment == "qa" && terraform.workspace != "default" ? "-${terraform.workspace}" : ""
  service_name        = "${var.service_name}${local.service_name_suffix}"

  full_name = "${var.project_name}-${var.environment}-${local.service_name}"

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
