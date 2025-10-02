/**
 * ECS Cluster module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "cluster_name" {
  description = "Name of the ECS cluster (defaults to {project}-{env}-ecs)"
  type        = string
  default     = ""
}

variable "capacity_providers" {
  description = "List of capacity providers to use for the cluster"
  type        = list(string)
  default     = ["FARGATE", "FARGATE_SPOT"]
}

variable "default_capacity_provider_strategy" {
  description = "Default capacity provider strategy for the cluster"
  type = list(object({
    capacity_provider = string
    weight            = number
    base              = number
  }))
  default = [
    {
      capacity_provider = "FARGATE"
      weight            = 1
      base              = 1
    },
    {
      capacity_provider = "FARGATE_SPOT"
      weight            = 4
      base              = 0
    }
  ]
}

variable "container_insights" {
  description = "Whether to enable CloudWatch Container Insights"
  type        = bool
  default     = true
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
