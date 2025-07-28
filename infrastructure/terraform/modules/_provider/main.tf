/**
 * AWS Provider Configuration with LocalStack support
 * This module can be included in other modules to provide LocalStack compatibility
 */

variable "aws_endpoint" {
  description = "Custom AWS endpoint URL (for LocalStack)"
  type        = string
  default     = null
}

variable "skip_credentials_validation" {
  description = "Skip credentials validation (for LocalStack)"
  type        = bool
  default     = false
}

variable "skip_metadata_api_check" {
  description = "Skip metadata API check (for LocalStack)"
  type        = bool
  default     = false
}

variable "skip_region_validation" {
  description = "Skip region validation (for LocalStack)"
  type        = bool
  default     = false
}

variable "aws_region" {
  description = "AWS region to use"
  type        = string
  default     = "us-east-1"
}

# Determine if we're using LocalStack based on whether aws_endpoint is set
locals {
  using_localstack = var.aws_endpoint != null
}

# AWS Provider configuration
provider "aws" {
  region     = var.aws_region
  access_key = local.using_localstack ? "test" : null
  secret_key = local.using_localstack ? "test" : null

  # LocalStack specific settings
  s3_use_path_style           = local.using_localstack
  skip_credentials_validation = local.using_localstack ? true : var.skip_credentials_validation
  skip_metadata_api_check     = local.using_localstack ? true : var.skip_metadata_api_check
  skip_region_validation      = local.using_localstack ? true : var.skip_region_validation

  # Only set these if we're using LocalStack
  dynamic "endpoints" {
    for_each = local.using_localstack ? [1] : []
    content {
      apigateway     = var.aws_endpoint
      cloudformation = var.aws_endpoint
      cloudwatch     = var.aws_endpoint
      dynamodb       = var.aws_endpoint
      ec2            = var.aws_endpoint
      ecs            = var.aws_endpoint
      eks            = var.aws_endpoint
      iam            = var.aws_endpoint
      lambda         = var.aws_endpoint
      route53        = var.aws_endpoint
      s3             = var.aws_endpoint
      secretsmanager = var.aws_endpoint
      ses            = var.aws_endpoint
      sns            = var.aws_endpoint
      sqs            = var.aws_endpoint
      ssm            = var.aws_endpoint
      rds            = var.aws_endpoint
    }
  }
}
