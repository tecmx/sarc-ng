/**
 * Example usage of the database module
 */

provider "aws" {
  region = "us-west-2"
}

# Use data source to get the VPC and subnets
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Example of standard RDS instance
module "mysql_db" {
  source = "../.."

  project_name = "sarc"
  environment  = "dev"

  vpc_id     = data.aws_vpc.default.id
  subnet_ids = data.aws_subnets.private.ids

  # Database configuration
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = "db.t3.medium"

  # Security
  allowed_cidr_blocks = ["10.0.0.0/16"]
  deletion_protection = false

  additional_tags = {
    Example = "true"
  }
}

# Example of Aurora cluster
module "aurora_db" {
  source = "../.."

  project_name = "sarc"
  environment  = "staging"

  vpc_id     = data.aws_vpc.default.id
  subnet_ids = data.aws_subnets.private.ids

  # Use Aurora instead of standard RDS
  is_aurora      = true
  engine         = "aurora-mysql"
  engine_version = "8.0.mysql_aurora.3.03.1"
  instance_class = "db.r5.large"

  # Security
  allowed_cidr_blocks        = ["10.0.0.0/16"]
  allowed_security_group_ids = ["sg-12345678"]
  deletion_protection        = false

  additional_tags = {
    Example = "true"
  }
} 
