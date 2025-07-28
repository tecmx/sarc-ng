output "lambda_function_arn" {
  description = "The ARN of the Lambda function"
  value       = aws_lambda_function.this.arn
}

output "lambda_function_name" {
  description = "The name of the Lambda function"
  value       = aws_lambda_function.this.function_name
}

output "lambda_function_invoke_arn" {
  description = "The invoke ARN of the Lambda function"
  value       = aws_lambda_function.this.invoke_arn
}

output "lambda_function_version" {
  description = "The version of the Lambda function"
  value       = aws_lambda_function.this.version
}

output "lambda_function_last_modified" {
  description = "The date the Lambda function was last modified"
  value       = aws_lambda_function.this.last_modified
}

output "lambda_role_arn" {
  description = "The ARN of the IAM role created for the Lambda function"
  value       = var.iam_role_arn != null ? var.iam_role_arn : try(aws_iam_role.lambda[0].arn, null)
}

output "lambda_role_name" {
  description = "The name of the IAM role created for the Lambda function"
  value       = var.iam_role_arn != null ? null : try(aws_iam_role.lambda[0].name, null)
}

output "event_source_mapping_id" {
  description = "The ID of the Lambda event source mapping"
  value       = aws_lambda_event_source_mapping.this.id
}

output "event_source_mapping_uuid" {
  description = "The UUID of the Lambda event source mapping"
  value       = aws_lambda_event_source_mapping.this.uuid
}

output "event_source_mapping_function_arn" {
  description = "The ARN of the Lambda function the event source mapping is sending events to"
  value       = aws_lambda_event_source_mapping.this.function_arn
}

output "dlq_arn" {
  description = "The ARN of the dead-letter queue"
  value       = var.create_dlq ? aws_sqs_queue.dlq[0].arn : null
}

output "dlq_id" {
  description = "The ID of the dead-letter queue"
  value       = var.create_dlq ? aws_sqs_queue.dlq[0].id : null
}

output "dlq_url" {
  description = "The URL of the dead-letter queue"
  value       = var.create_dlq ? aws_sqs_queue.dlq[0].url : null
} 
