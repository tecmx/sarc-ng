/**
 * DNS Record module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "zone_id" {
  description = "ID of the hosted zone where the record will be created"
  type        = string
}

variable "zone_name" {
  description = "Name of the hosted zone (e.g., example.com)"
  type        = string
}

variable "record_name" {
  description = "Name of the record without the domain (e.g., 'api' will create 'api.example.com')"
  type        = string
}

variable "record_type" {
  description = "Type of DNS record (A, AAAA, CNAME, MX, TXT, etc.)"
  type        = string
  default     = "A"
}

variable "record_ttl" {
  description = "TTL for the record in seconds"
  type        = number
  default     = 60
}

variable "records" {
  description = "List of record values (not used for alias records)"
  type        = list(string)
  default     = []
}

variable "alias" {
  description = "Alias record configuration"
  type = object({
    name                   = string
    zone_id                = string
    evaluate_target_health = bool
  })
  default = null
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
