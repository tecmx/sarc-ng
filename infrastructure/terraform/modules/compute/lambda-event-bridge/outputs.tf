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

output "event_rule_id" {
  description = "The ID of the EventBridge rule"
  value       = aws_cloudwatch_event_rule.this.id
}

output "event_rule_arn" {
  description = "The ARN of the EventBridge rule"
  value       = aws_cloudwatch_event_rule.this.arn
}

output "event_target_id" {
  description = "The ID of the EventBridge target"
  value       = aws_cloudwatch_event_target.this.id
} 
