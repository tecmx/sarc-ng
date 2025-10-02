/**
 * Database module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

  # Database name - sanitized to remove hyphens
  db_name = replace("${var.project_name}_${var.environment}", "-", "_")

  # Common tags for all resources
  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      Terraform   = "true"
    },
    var.additional_tags
  )
}
