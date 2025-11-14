# Dev environment specific configuration for LocalStack

# Define locals to make the configuration more DRY (replicated from root)
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

# Note: Using local state for LocalStack to avoid authentication issues

# Override provider configuration for LocalStack
generate "provider" {
  path = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents = <<EOF
provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = "us-east-1"
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    apigateway     = "http://localhost:4566"
    cloudformation = "http://localhost:4566"
    dynamodb       = "http://localhost:4566"
    ec2            = "http://localhost:4566"
    iam            = "http://localhost:4566"
    kinesis        = "http://localhost:4566"
    lambda         = "http://localhost:4566"
    rds            = "http://localhost:4566"
    route53        = "http://localhost:4566"
    s3             = "http://localhost:4566"
    ses            = "http://localhost:4566"
    sns            = "http://localhost:4566"
    sqs            = "http://localhost:4566"
    ssm            = "http://localhost:4566"
    sts            = "http://localhost:4566"
  }
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
  is_local_development = true
  
  # Pass the parsed values as inputs to all child modules
  aws_region = local.aws_region
  account_id = local.account_id
  environment = local.environment
} 