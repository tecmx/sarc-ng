# Network module configuration for dev environment

include {
  path = "${get_path_to_repo_root()}/infrastructure/terraform/live/accounts/dev/terragrunt.hcl"
  expose = true
  merge_strategy = "deep"
}

# Use the network module
terraform {
  source = "${get_path_to_repo_root()}/infrastructure/terraform/modules/network"
  
  before_hook "validate_vpc_cidr" {
    commands = ["plan", "apply"]
    execute  = [
      "bash", "-c", 
      "echo Validating VPC CIDR block ${include.locals.env_vars.locals.vpc_cidr}"
    ]
  }
}

# Set dependencies explicitly
dependencies {
  paths = []
}

# Network-specific inputs
inputs = {
  # Basic network configuration from env.hcl
  vpc_cidr = include.locals.env_vars.locals.vpc_cidr
  availability_zones = include.locals.env_vars.locals.availability_zones
  public_subnet_cidrs = include.locals.env_vars.locals.public_subnet_cidrs
  private_subnet_cidrs = include.locals.env_vars.locals.private_subnet_cidrs
  database_subnet_cidrs = include.locals.env_vars.locals.database_subnet_cidrs
  
  # VPC features configuration
  enable_vpn_gateway   = false
  
  # Disable database subnet group for LocalStack (RDS not supported in Community edition)
  create_database_subnet_group = false
  
  # Disable VPC flow logs for LocalStack (CloudWatch Logs not needed for local dev)
  enable_flow_log = false
  
  # Additional subnet tags for Kubernetes integration
  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }
  
  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }
  
  # Additional tags
  additional_tags = include.locals.env_vars.locals.additional_tags
} 