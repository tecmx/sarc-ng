/**
 * Database module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC where database will be deployed"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for the DB subnet group"
  type        = list(string)
}

variable "engine" {
  description = "Database engine (mysql, postgres, aurora-mysql, aurora-postgresql)"
  type        = string
  default     = "mysql"
}

variable "engine_version" {
  description = "Database engine version"
  type        = string
  default     = "8.0"
}

variable "instance_class" {
  description = "Instance type for the RDS instances"
  type        = string
  default     = "db.t3.medium"
}

variable "allocated_storage" {
  description = "Allocated storage in GB (not applicable for Aurora)"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Max allocated storage for autoscaling (not applicable for Aurora)"
  type        = number
  default     = 100
}

variable "is_aurora" {
  description = "Whether to create an Aurora cluster instead of a standard RDS instance"
  type        = bool
  default     = false
}

variable "multi_az" {
  description = "Whether to deploy a multi-AZ RDS instance (not applicable for Aurora)"
  type        = bool
  default     = false
}

variable "backup_retention_period" {
  description = "Days to retain backups"
  type        = number
  default     = 7
}

variable "deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = true
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "master_password" {
  description = "Master password for the database (if not provided, a random one will be generated)"
  type        = string
  sensitive   = true
  default     = null
}

variable "port" {
  description = "Database port (3306 for MySQL, 5432 for PostgreSQL)"
  type        = number
  default     = 3306
}

variable "allowed_cidr_blocks" {
  description = "List of CIDR blocks that are allowed to access the database"
  type        = list(string)
  default     = []
}

variable "allowed_security_group_ids" {
  description = "List of security group IDs that are allowed to access the database"
  type        = list(string)
  default     = []
}
