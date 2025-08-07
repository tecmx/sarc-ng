locals {
  api_name        = var.api_name != null ? var.api_name : "${var.name}-api"
  api_description = var.api_description != null ? var.api_description : "API Gateway for ${var.name} Lambda function"

  lambda_source = var.source_path != null ? {
    path = var.source_path
    } : var.source_bucket != null && var.source_key != null ? {
    s3_bucket = var.source_bucket
    s3_key    = var.source_key
  } : null

  lambda_role_name = "${var.name}-lambda-role"

  default_cors_configuration = {
    allow_origins     = ["*"]
    allow_methods     = ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"]
    allow_headers     = ["Content-Type", "Authorization", "X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token"]
    expose_headers    = ["Content-Type", "X-Amz-Date", "X-Api-Key"]
    allow_credentials = false
    max_age           = 300
  }

  cors_configuration = var.cors_configuration != null ? var.cors_configuration : local.default_cors_configuration

  tags = merge(
    var.tags,
    {
      Name      = var.name
      ManagedBy = "terraform"
    }
  )
} 
