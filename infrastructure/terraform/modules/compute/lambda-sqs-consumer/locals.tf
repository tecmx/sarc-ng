locals {
  lambda_source = var.source_path != null ? {
    path = var.source_path
    } : var.source_bucket != null && var.source_key != null ? {
    s3_bucket = var.source_bucket
    s3_key    = var.source_key
  } : null

  lambda_role_name = "${var.name}-lambda-role"
  dlq_name         = var.dlq_name != null ? var.dlq_name : "${var.name}-dlq"

  # Extract SQS queue name from ARN
  sqs_queue_name = element(split(":", var.sqs_queue_arn), length(split(":", var.sqs_queue_arn)) - 1)

  tags = merge(
    var.tags,
    {
      Name      = var.name
      ManagedBy = "terraform"
    }
  )
} 
