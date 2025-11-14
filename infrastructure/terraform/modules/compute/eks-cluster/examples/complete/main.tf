/**
 * Example usage of the eks-cluster module
 */

provider "aws" {
  region = "us-east-1"
}

# Assume VPC exists
data "aws_vpc" "main" {
  tags = {
    Name = "sarc-ng-dev"
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

# Create EKS cluster with managed node groups
module "eks" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  cluster_name    = "sarc-ng-dev-eks"
  cluster_version = "1.28"

  # Networking
  vpc_id          = data.aws_vpc.main.id
  private_subnets = data.aws_subnets.private.ids

  # Node groups with SPOT instances for cost savings
  node_groups = {
    general = {
      name           = "general"
      min_size       = 2
      max_size       = 4
      desired_size   = 2
      instance_types = ["t3.medium"]
      capacity_type  = "SPOT"
    }
  }

  # Public access for development
  cluster_endpoint_public_access = true

  additional_tags = {
    Example = "complete"
  }
}

