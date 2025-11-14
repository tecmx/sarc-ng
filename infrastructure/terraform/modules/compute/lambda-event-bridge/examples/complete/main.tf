/**
 * Example usage of the lambda-event-bridge module
 */

provider "aws" {
  region = "us-east-1"
}

# Lambda function triggered by EventBridge rule
module "scheduled_task" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"
  name         = "daily-report"

  description = "Generate daily reports"
  handler     = "index.handler"
  runtime     = "python3.11"
  memory_size = 512
  timeout     = 300

  # Source code
  source_path = "${path.module}/function.zip"

  # EventBridge schedule (daily at 2 AM)
  schedule_expression = "cron(0 2 * * ? *)"
  rule_name           = "daily-report-trigger"

  # Enable X-Ray
  enable_xray_tracing = true

  additional_tags = {
    Example = "complete"
    Type    = "scheduled"
  }
}

