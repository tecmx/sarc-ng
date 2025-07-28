# Environment-specific variables
locals {
  environment = "dev"
  
  # Network configuration
  vpc_cidr                 = "10.0.0.0/16"
  availability_zones       = ["us-east-1a", "us-east-1b"]
  public_subnet_cidrs      = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidrs     = ["10.0.10.0/24", "10.0.11.0/24"]
  database_subnet_cidrs    = ["10.0.20.0/24", "10.0.21.0/24"]
  single_nat_gateway       = true  # Cost optimization for dev environment
  
  # Additional tags
  additional_tags = {
    CostCenter = "DevOps"
    Team       = "Engineering"
  }
} 