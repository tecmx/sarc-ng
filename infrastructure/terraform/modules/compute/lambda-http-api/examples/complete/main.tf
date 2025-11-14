/**
 * Example usage of the lambda-http-api module
 */

provider "aws" {
  region = "us-east-1"
}

# Lambda function with HTTP API Gateway
module "api_lambda" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"
  name         = "hello-api"

  description = "Example Lambda function with HTTP API"
  handler     = "index.handler"
  runtime     = "python3.11"
  memory_size = 256
  timeout     = 10

  # Source code (create a simple handler)
  source_path = "${path.module}/function.zip"

  # API Gateway configuration
  create_api_gateway = true
  api_name           = "hello-api"
  api_stage_name     = "v1"

  # CORS configuration
  cors_configuration = {
    allow_origins     = ["https://app.example.com"]
    allow_methods     = ["GET", "POST"]
    allow_headers     = ["Content-Type", "Authorization"]
    expose_headers    = []
    allow_credentials = false
    max_age           = 300
  }

  # Logging
  log_retention_in_days = 7

  # X-Ray tracing
  enable_xray_tracing = true

  additional_tags = {
    Example = "complete"
  }
}

