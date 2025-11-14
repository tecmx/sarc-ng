/**
 * Example usage of the ecs-alb-service module
 */

provider "aws" {
  region = "us-east-1"
}

# Assume VPC and ECS cluster exist
data "aws_vpc" "main" {
  tags = {
    Name = "sarc-ng-dev"
  }
}

data "aws_subnets" "public" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
  tags = {
    Tier = "public"
  }
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.main.id]
  }
  tags = {
    Tier = "private"
  }
}

# Deploy ECS service with ALB
module "api_service" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"
  service_name = "api"

  # Cluster configuration
  cluster_id   = "arn:aws:ecs:us-east-1:123456789012:cluster/sarc-ng-dev-ecs"
  cluster_name = "sarc-ng-dev-ecs"

  # Networking
  vpc_id          = data.aws_vpc.main.id
  public_subnets  = data.aws_subnets.public.ids
  private_subnets = data.aws_subnets.private.ids

  # Container configuration
  container_image = "nginx:latest"
  container_port  = 80
  cpu             = 256
  memory          = 512

  # Scaling
  desired_count = 2
  min_capacity  = 1
  max_capacity  = 4

  # Health check
  health_check_path = "/health"

  additional_tags = {
    Example = "complete"
  }
}

