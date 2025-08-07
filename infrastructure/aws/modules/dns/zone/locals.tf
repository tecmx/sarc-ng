/**
 * DNS Zone module - Local variables
 */

locals {
  name = "${var.project_name}-${var.environment}"

  # Full list of domains for the certificate
  # Always include wildcard for the zone
  cert_domains = distinct(concat(["*.${var.zone_name}", var.zone_name], var.subject_alternative_names))

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    },
    var.additional_tags
  )
} 
