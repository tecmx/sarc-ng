/**
 * Database Schema module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"
  # If user_name is not specified, use schema_name
  user = var.user_name == "" ? var.schema_name : var.user_name

  # For QA environments with workspaces, add workspace suffix
  schema_suffix = var.environment == "qa" && terraform.workspace != "default" ? "_${terraform.workspace}" : ""
  schema_name   = "${var.schema_name}${local.schema_suffix}"
  user_name     = "${local.user}${local.schema_suffix}"

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Service     = var.schema_name
    },
    var.additional_tags
  )
} 
