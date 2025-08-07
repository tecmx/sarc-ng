/**
 * EKS Cluster module - Input variables
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
  description = "ID of the VPC where the cluster will be deployed"
  type        = string
}

variable "private_subnets" {
  description = "List of private subnet IDs where the cluster and nodes will be deployed"
  type        = list(string)
}

variable "public_subnets" {
  description = "List of public subnet IDs where the public load balancers will be deployed"
  type        = list(string)
}

variable "cluster_name" {
  description = "Name of the EKS cluster (defaults to {project}-{env}-eks)"
  type        = string
  default     = ""
}

variable "cluster_version" {
  description = "Kubernetes version to use for the cluster"
  type        = string
  default     = "1.28"
}

variable "cluster_endpoint_public_access" {
  description = "Whether the cluster's API server is accessible publicly"
  type        = bool
  default     = true
}

variable "cluster_endpoint_private_access" {
  description = "Whether the cluster's API server is accessible privately"
  type        = bool
  default     = true
}

variable "cluster_endpoint_public_access_cidrs" {
  description = "List of CIDR blocks that can access the cluster API server publicly"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "node_group_defaults" {
  description = "Default configuration for managed node groups"
  type        = any
  default = {
    disk_size      = 50
    instance_types = ["t3.medium"]
    capacity_type  = "SPOT"
    min_size       = 1
    max_size       = 3
    desired_size   = 1
  }
}

variable "node_groups" {
  description = "Map of managed node group configurations"
  type        = any
  default     = {}
}

variable "aws_auth_roles" {
  description = "List of IAM roles to add to the aws-auth ConfigMap"
  type = list(object({
    rolearn  = string
    username = string
    groups   = list(string)
  }))
  default = []
}

variable "aws_auth_users" {
  description = "List of IAM users to add to the aws-auth ConfigMap"
  type = list(object({
    userarn  = string
    username = string
    groups   = list(string)
  }))
  default = []
}

variable "additional_tags" {
  description = "Additional tags to add to all resources"
  type        = map(string)
  default     = {}
} 
