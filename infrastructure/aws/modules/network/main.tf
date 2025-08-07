/**
 * Network module - Creates a VPC with subnets and other networking components
 */

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = ">= 5.0"

  name = "${var.project_name}-${var.environment}"
  cidr = var.vpc_cidr

  azs              = var.availability_zones
  private_subnets  = var.private_subnet_cidrs
  public_subnets   = var.public_subnet_cidrs
  database_subnets = var.database_subnet_cidrs

  # Public network access configuration
  create_igw              = true
  enable_nat_gateway      = true
  single_nat_gateway      = var.environment != "prod"
  one_nat_gateway_per_az  = var.environment == "prod"
  enable_dns_hostnames    = true
  enable_dns_support      = true
  map_public_ip_on_launch = true

  # VPC endpoints
  enable_vpn_gateway = var.enable_vpn_gateway

  # VPC Flow Logs configuration
  enable_flow_log                      = var.enable_flow_log
  create_flow_log_cloudwatch_iam_role  = var.enable_flow_log
  create_flow_log_cloudwatch_log_group = var.enable_flow_log

  # Database configuration
  create_database_subnet_group       = var.create_database_subnet_group
  create_database_subnet_route_table = true

  # Tags
  tags = local.tags

  vpc_tags = {
    Name = "${var.project_name}-${var.environment}-vpc"
  }

  igw_tags = {
    Name = "${var.project_name}-${var.environment}-igw"
  }

  nat_gateway_tags = {
    Name = "${var.project_name}-${var.environment}-nat"
  }

  public_subnet_tags = merge(
    {
      Name                     = "${var.project_name}-${var.environment}-public"
      "kubernetes.io/role/elb" = "1"
    },
    var.public_subnet_tags
  )

  private_subnet_tags = merge(
    {
      Name                              = "${var.project_name}-${var.environment}-private"
      "kubernetes.io/role/internal-elb" = "1"
    },
    var.private_subnet_tags
  )

  database_subnet_tags = merge(
    {
      Name = "${var.project_name}-${var.environment}-db"
    },
    var.database_subnet_tags
  )
}
