/**
 * ECS NLB Service module - Input variables
 */

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, qa, staging, prod)"
  type        = string
}

variable "service_name" {
  description = "Name of the service"
  type        = string
}

variable "cluster_arn" {
  description = "ARN of the ECS cluster"
  type        = string
}

variable "vpc_id" {
  description = "ID of the VPC where the service will be deployed"
  type        = string
}

variable "private_subnets" {
  description = "List of private subnet IDs where the service will be deployed"
  type        = list(string)
}

variable "public_subnets" {
  description = "List of public subnet IDs where the NLB will be deployed"
  type        = list(string)
}

variable "container_image" {
  description = "Docker image to run (repo/image:tag)"
  type        = string
}

variable "container_port" {
  description = "Port the container exposes"
  type        = number
  default     = 8080
}

variable "listener_port" {
  description = "Port the NLB will listen on"
  type        = number
  default     = 80
}

variable "cpu" {
  description = "CPU units for the task (1024 = 1 vCPU)"
  type        = number
  default     = 512
}

variable "memory" {
  description = "Memory for the task in MB"
  type        = number
  default     = 1024
}

variable "desired_count" {
  description = "Desired number of tasks"
  type        = number
  default     = 1
}

variable "health_check_path" {
  description = "Path for NLB health check"
  type        = string
  default     = "/health"
}

variable "health_check_port" {
  description = "Port for NLB health check"
  type        = number
  default     = 0 # 0 means same as traffic port
}

variable "container_environment" {
  description = "Environment variables for the container"
  type        = map(string)
  default     = {}
  sensitive   = true
}

variable "container_secrets" {
  description = "Secrets (from SSM or Secrets Manager) to inject into the container"
  type = list(object({
    name      = string
    valueFrom = string
  }))
  default = []
}

variable "assign_public_ip" {
  description = "Whether to assign public IP addresses to the task"
  type        = bool
  default     = false
}

variable "enable_execution_role" {
  description = "Whether to create an execution role for the task"
  type        = bool
  default     = true
}

variable "execution_role_arn" {
  description = "Existing execution role ARN to use (if not creating one)"
  type        = string
  default     = null
}

variable "platform_version" {
  description = "The platform version to run on Fargate"
  type        = string
  default     = "LATEST"
}

variable "capacity_provider_strategy" {
  description = "Capacity provider strategy for the service"
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
    }
  ]
}

variable "deployment_maximum_percent" {
  description = "Maximum percentage of tasks during deployment"
  type        = number
  default     = 200
}

variable "deployment_minimum_healthy_percent" {
  description = "Minimum percentage of tasks that must remain healthy during deployment"
  type        = number
  default     = 100
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
