/**
 * DNS Zone module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "zone_name" {
  description = "Domain name for the hosted zone (e.g., example.com or dev.example.com)"
  type        = string
}

variable "create_public_zone" {
  description = "Whether to create a public hosted zone (true) or private hosted zone (false)"
  type        = bool
  default     = true
}

variable "vpc_id" {
  description = "VPC ID for private hosted zone (required if create_public_zone is false)"
  type        = string
  default     = null
}

variable "create_acm_certificate" {
  description = "Whether to create an ACM certificate for the zone"
  type        = bool
  default     = true
}

variable "subject_alternative_names" {
  description = "Additional domain names for the ACM certificate"
  type        = list(string)
  default     = []
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
