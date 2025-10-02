/**
 * Database Schema module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "admin_secret_arn" {
  description = "ARN of the secret containing admin credentials"
  type        = string
}

variable "schema_name" {
  description = "Name of the schema to create"
  type        = string
}

variable "user_name" {
  description = "Name of the user to create"
  type        = string
  default     = "" # If empty, will use schema_name
}

variable "host" {
  description = "Database host"
  type        = string
}

variable "port" {
  description = "Database port"
  type        = number
}

variable "engine" {
  description = "Database engine (mysql or postgres)"
  type        = string
  default     = "mysql"
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
