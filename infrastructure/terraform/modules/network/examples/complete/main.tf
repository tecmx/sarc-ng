/**
 * Example usage of the network module
 */

provider "aws" {
  region = "us-east-1"
}

# Basic VPC with public, private, and database subnets
module "vpc" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  vpc_cidr           = "10.0.0.0/16"
  availability_zones = ["us-east-1a", "us-east-1b"]

  public_subnet_cidrs   = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidrs  = ["10.0.10.0/24", "10.0.11.0/24"]
  database_subnet_cidrs = ["10.0.20.0/24", "10.0.21.0/24"]

  # VPC features
  enable_vpn_gateway = false
  enable_flow_log    = true

  # Subnet tagging for EKS
  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
  }

  additional_tags = {
    Example = "complete"
  }
}

