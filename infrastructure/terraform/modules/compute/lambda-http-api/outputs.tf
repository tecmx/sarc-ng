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

output "api_gateway_id" {
  description = "The ID of the API Gateway"
  value       = var.create_api_gateway ? aws_apigatewayv2_api.this[0].id : null
}

output "api_gateway_arn" {
  description = "The ARN of the API Gateway"
  value       = var.create_api_gateway ? aws_apigatewayv2_api.this[0].arn : null
}

output "api_gateway_execution_arn" {
  description = "The execution ARN of the API Gateway"
  value       = var.create_api_gateway ? aws_apigatewayv2_api.this[0].execution_arn : null
}

output "api_gateway_endpoint" {
  description = "The endpoint URL of the API Gateway"
  value       = var.create_api_gateway ? aws_apigatewayv2_stage.this[0].invoke_url : null
}

output "api_gateway_stage_id" {
  description = "The ID of the API Gateway stage"
  value       = var.create_api_gateway ? aws_apigatewayv2_stage.this[0].id : null
}

output "api_gateway_stage_arn" {
  description = "The ARN of the API Gateway stage"
  value       = var.create_api_gateway ? aws_apigatewayv2_stage.this[0].arn : null
}

output "api_gateway_domain_name" {
  description = "The custom domain name of the API Gateway"
  value       = var.create_api_gateway && var.create_custom_domain && var.api_domain_name != null ? aws_apigatewayv2_domain_name.this[0].domain_name : null
}

output "api_gateway_domain_name_target" {
  description = "The target domain name of the API Gateway custom domain"
  value       = var.create_api_gateway && var.create_custom_domain && var.api_domain_name != null ? try(aws_apigatewayv2_domain_name.this[0].domain_name_configuration[0].target_domain_name, null) : null
}

output "api_gateway_domain_name_hosted_zone_id" {
  description = "The hosted zone ID of the API Gateway custom domain"
  value       = var.create_api_gateway && var.create_custom_domain && var.api_domain_name != null ? try(aws_apigatewayv2_domain_name.this[0].domain_name_configuration[0].hosted_zone_id, null) : null
} 
