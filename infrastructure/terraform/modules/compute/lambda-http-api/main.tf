/**
 * Lambda HTTP API module - Lambda function resources
 */

#################
# Lambda Function
#################

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = var.log_retention_in_days
  tags              = local.tags
}

resource "aws_lambda_function" "this" {
  function_name = var.name
  description   = var.description
  role          = var.iam_role_arn != null ? var.iam_role_arn : aws_iam_role.lambda[0].arn
  handler       = var.handler
  runtime       = var.runtime
  memory_size   = var.memory_size
  timeout       = var.timeout
  publish       = var.publish

  dynamic "vpc_config" {
    for_each = var.vpc_config != null ? [var.vpc_config] : []
    content {
      subnet_ids         = vpc_config.value.subnet_ids
      security_group_ids = vpc_config.value.security_group_ids
    }
  }

  dynamic "environment" {
    for_each = length(var.environment_variables) > 0 ? [var.environment_variables] : []
    content {
      variables = environment.value
    }
  }

  dynamic "dead_letter_config" {
    for_each = var.destination_on_failure != null ? [var.destination_on_failure] : []
    content {
      target_arn = dead_letter_config.value
    }
  }

  tracing_config {
    mode = var.enable_xray_tracing ? "Active" : "PassThrough"
  }

  reserved_concurrent_executions = var.reserved_concurrent_executions

  # Source code
  dynamic "s3_key" {
    for_each = try(local.lambda_source.s3_key, null) != null ? [local.lambda_source.s3_key] : []
    content {
      s3_bucket = local.lambda_source.s3_bucket
      s3_key    = local.lambda_source.s3_key
    }
  }

  dynamic "filename" {
    for_each = try(local.lambda_source.path, null) != null ? [local.lambda_source.path] : []
    content {
      filename = local.lambda_source.path
    }
  }

  # Use source_code_hash if available
  source_code_hash = try(filebase64sha256(local.lambda_source.path), null)

  depends_on = [
    aws_cloudwatch_log_group.lambda,
    aws_iam_role_policy_attachment.logs,
    aws_iam_role_policy_attachment.vpc_access,
    aws_iam_role_policy_attachment.xray,
    aws_iam_role_policy_attachment.custom
  ]

  tags = local.tags
}

resource "aws_lambda_function_event_invoke_config" "this" {
  count = var.create_async_event_config ? 1 : 0

  function_name                = aws_lambda_function.this.function_name
  qualifier                    = aws_lambda_function.this.version
  maximum_retry_attempts       = var.maximum_retry_attempts
  maximum_event_age_in_seconds = var.maximum_event_age_in_seconds

  dynamic "destination_config" {
    for_each = var.destination_on_failure != null || var.destination_on_success != null ? [true] : []
    content {
      dynamic "on_failure" {
        for_each = var.destination_on_failure != null ? [var.destination_on_failure] : []
        content {
          destination = on_failure.value
        }
      }

      dynamic "on_success" {
        for_each = var.destination_on_success != null ? [var.destination_on_success] : []
        content {
          destination = on_success.value
        }
      }
    }
  }
}
