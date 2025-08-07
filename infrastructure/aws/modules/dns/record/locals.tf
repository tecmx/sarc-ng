/**
 * DNS Record module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

  # For QA environments with workspaces, add workspace prefix
  prefix = var.environment == "qa" && terraform.workspace != "default" ? "${terraform.workspace}." : ""

  # Full record name
  full_record_name = "${local.prefix}${var.record_name}.${var.zone_name}"

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Service     = var.record_name
    },
    var.additional_tags
  )
} 
