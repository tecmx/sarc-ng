# Root terragrunt.hcl configuration for all environments

# Define locals to make the configuration more DRY
locals {
  # Parse account and region information
  account_vars = read_terragrunt_config(find_in_parent_folders("account.hcl", "empty.hcl"))
  region_vars = read_terragrunt_config(find_in_parent_folders("region.hcl", "empty.hcl"))
  env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl", "empty.hcl"))
  
  # Extract commonly used variables
  account_name = local.account_vars.locals.account_name
  account_id = try(local.account_vars.locals.aws_account_id, "")
  aws_region = try(local.region_vars.locals.aws_region, "us-east-1")
  environment = try(local.env_vars.locals.environment, "")
}

# Configure Terragrunt to use remote state in S3 with DynamoDB locking
remote_state {
  backend = "s3"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
  config = {
    bucket         = "sarc-ng-terraform-state-${local.account_name}-${local.aws_region}"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = local.aws_region
    encrypt        = true
    dynamodb_table = "sarc-ng-terraform-locks-${local.account_name}"
  }
}

# Generate provider configuration
generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "aws" {
  region = var.region

  # Skip credentials validation and account ID retrieval for local development
  skip_credentials_validation = var.is_local_development
  skip_requesting_account_id  = var.is_local_development

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
      Workspace   = terraform.workspace
    }
  }
}

variable "is_local_development" {
  description = "Whether we're running in a local development environment"
  type        = bool
  default     = false
}
EOF
}

# Generate required variables that will be shared across all modules
generate "common_variables" {
  path      = "common_variables.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "region" {
  description = "AWS region to deploy resources to"
  type        = string
}
EOF
}

# Standard validation and preparation hooks
terraform {
  before_hook "before_hook" {
    commands     = ["apply", "plan"]
    execute      = ["echo", "Running Terraform on ${local.environment} environment in ${local.aws_region}"]
  }
  
  after_hook "after_hook" {
    commands     = ["apply"]
    execute      = ["echo", "Terraform apply completed successfully!"]
    run_on_error = false
  }
}

# Common inputs for all modules
inputs = {
  project_name = "sarc-ng"
  is_local_development = false
  
  # Pass the parsed values as inputs to all child modules
  aws_region = local.aws_region
  account_id = local.account_id
  environment = local.environment
} 