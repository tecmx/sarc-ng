/**
 * Network module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

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
