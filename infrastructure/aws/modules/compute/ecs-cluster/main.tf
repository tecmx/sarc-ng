/**
 * ECS Cluster module - Creates an ECS cluster for running Fargate services
 */

# ECS Cluster
module "ecs" {
  source  = "terraform-aws-modules/ecs/aws"
  version = "~> 5.0"

  cluster_name = local.cluster_name

  # Capacity provider strategy
  fargate_capacity_providers = {
    FARGATE = {
      default_capacity_provider_strategy = {
        weight = lookup(
          { for s in var.default_capacity_provider_strategy : s.capacity_provider => s },
          "FARGATE",
          { weight = 1 }
        ).weight
        base = lookup(
          { for s in var.default_capacity_provider_strategy : s.capacity_provider => s },
          "FARGATE",
          { base = 0 }
        ).base
      }
    }
    FARGATE_SPOT = {
      default_capacity_provider_strategy = {
        weight = lookup(
          { for s in var.default_capacity_provider_strategy : s.capacity_provider => s },
          "FARGATE_SPOT",
          { weight = 0 }
        ).weight
        base = lookup(
          { for s in var.default_capacity_provider_strategy : s.capacity_provider => s },
          "FARGATE_SPOT",
          { base = 0 }
        ).base
      }
    }
  }

  # CloudWatch Container Insights
  cluster_settings = {
    name  = "containerInsights"
    value = var.container_insights ? "enabled" : "disabled"
  }

  tags = local.tags
}
