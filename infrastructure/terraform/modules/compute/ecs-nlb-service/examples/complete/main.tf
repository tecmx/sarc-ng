/**
 * Example usage of the ecs-nlb-service module
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

# Deploy ECS service with NLB
module "grpc_service" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"
  service_name = "grpc-api"

  # Cluster configuration
  cluster_id   = "arn:aws:ecs:us-east-1:123456789012:cluster/sarc-ng-dev-ecs"
  cluster_name = "sarc-ng-dev-ecs"

  # Networking
  vpc_id     = data.aws_vpc.main.id
  subnet_ids = data.aws_subnets.public.ids

  # Security - required CIDR blocks
  allowed_cidr_blocks = ["10.0.0.0/16"]

  # Container configuration
  container_image = "myapp/grpc-server:latest"
  container_port  = 50051
  cpu             = 512
  memory          = 1024

  # Scaling
  desired_count = 2

  additional_tags = {
    Example  = "complete"
    Protocol = "grpc"
  }
}

