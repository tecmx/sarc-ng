/**
 * Example usage of the lambda-sqs-consumer module
 */

provider "aws" {
  region = "us-east-1"
}

# Assume SQS queue exists
data "aws_sqs_queue" "tasks" {
  name = "sarc-ng-dev-tasks"
}

# Lambda function consuming SQS messages
module "task_processor" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"
  name         = "task-processor"

  description = "Process background tasks from SQS"
  handler     = "index.handler"
  runtime     = "python3.11"
  memory_size = 512
  timeout     = 60

  # Source code
  source_path = "${path.module}/function.zip"

  # SQS configuration
  sqs_queue_arn                  = data.aws_sqs_queue.tasks.arn
  batch_size                     = 10
  max_batching_window_in_seconds = 5

  # Dead Letter Queue
  create_dlq = true

  # Concurrent executions
  reserved_concurrent_executions = 5

  additional_tags = {
    Example = "complete"
    Type    = "sqs-consumer"
  }
}

