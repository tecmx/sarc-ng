/**
 * EKS Cluster module - Local variables
 */

locals {
  name         = "${var.project_name}-${var.environment}"
  cluster_name = var.cluster_name != "" ? var.cluster_name : "${local.name}-eks"

  # Default node groups if none specified
  default_node_groups = {
    general = {
      name           = "general"
      min_size       = 1
      max_size       = 3
      desired_size   = 1
      instance_types = ["t3.medium"]
      capacity_type  = "SPOT"
    }
  }

  # Use specified node groups or default
  node_groups = length(var.node_groups) > 0 ? var.node_groups : local.default_node_groups

  tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    },
    var.additional_tags
  )
} 
