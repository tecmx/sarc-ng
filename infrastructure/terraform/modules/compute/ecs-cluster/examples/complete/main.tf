/**
 * Example usage of the ecs-cluster module
 */

provider "aws" {
  region = "us-east-1"
}

# Create ECS cluster with Fargate capacity providers
module "ecs_cluster" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  cluster_name = "sarc-ng-dev-ecs"

  # Capacity providers
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  # Fargate cost optimization (80% SPOT, 20% on-demand)
  default_capacity_provider_strategy = [
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

  # Enable Container Insights for monitoring
  container_insights = true

  additional_tags = {
    Example = "complete"
  }
}

