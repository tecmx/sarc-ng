/**
 * EKS Namespace module - Input variables
 */

variable "namespace" {
  description = "The name of the Kubernetes namespace to create"
  type        = string
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
}

variable "labels" {
  description = "Additional labels to apply to the namespace"
  type        = map(string)
  default     = {}
}

variable "annotations" {
  description = "Annotations to apply to the namespace"
  type        = map(string)
  default     = {}
}

# Resource Quota
variable "create_resource_quota" {
  description = "Whether to create a resource quota for the namespace"
  type        = bool
  default     = false
}

variable "quota_requests_cpu" {
  description = "CPU request quota for the namespace"
  type        = string
  default     = "10"
}

variable "quota_requests_memory" {
  description = "Memory request quota for the namespace"
  type        = string
  default     = "20Gi"
}

variable "quota_limits_cpu" {
  description = "CPU limits quota for the namespace"
  type        = string
  default     = "20"
}

variable "quota_limits_memory" {
  description = "Memory limits quota for the namespace"
  type        = string
  default     = "40Gi"
}

variable "quota_pods" {
  description = "Maximum number of pods allowed in the namespace"
  type        = string
  default     = "50"
}

# RBAC
variable "create_namespace_admin_role" {
  description = "Whether to create a namespace admin role"
  type        = bool
  default     = false
}

variable "namespace_admins" {
  description = "List of users to assign as namespace admins"
  type        = list(string)
  default     = []
}

# Network Policies
variable "enable_network_policies" {
  description = "Whether to enable network policies for the namespace"
  type        = bool
  default     = false
}
