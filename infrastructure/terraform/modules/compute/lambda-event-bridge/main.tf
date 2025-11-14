###############
# IAM Resources
###############

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  count = var.iam_role_arn == null ? 1 : 0

  name               = local.lambda_role_name
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  tags               = local.tags
}

data "aws_iam_policy_document" "logs" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:*:*:*"]
  }
}

resource "aws_iam_policy" "logs" {
  count = var.iam_role_arn == null ? 1 : 0

  name   = "${var.name}-lambda-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  count = var.iam_role_arn == null ? 1 : 0

  role       = aws_iam_role.lambda[0].name
  policy_arn = aws_iam_policy.logs[0].arn
}

resource "aws_iam_role_policy_attachment" "vpc_access" {
  count = var.iam_role_arn == null && var.vpc_config != null ? 1 : 0

  role       = aws_iam_role.lambda[0].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_iam_role_policy_attachment" "xray" {
  count = var.iam_role_arn == null && var.enable_xray_tracing ? 1 : 0

  role       = aws_iam_role.lambda[0].name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_policy" "custom" {
  count = var.iam_role_arn == null && length(var.iam_policy_documents) > 0 ? length(var.iam_policy_documents) : 0

  name   = "${var.name}-lambda-custom-${count.index}"
  policy = var.iam_policy_documents[count.index]
}

resource "aws_iam_role_policy_attachment" "custom" {
  count = var.iam_role_arn == null && length(var.iam_policy_documents) > 0 ? length(var.iam_policy_documents) : 0

  role       = aws_iam_role.lambda[0].name
  policy_arn = aws_iam_policy.custom[count.index].arn
}

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

################
# EventBridge
################

resource "aws_cloudwatch_event_rule" "this" {
  name           = local.rule_name
  description    = local.rule_description
  event_pattern  = var.event_pattern
  event_bus_name = var.event_bus_name
  tags           = local.tags
}

resource "aws_cloudwatch_event_target" "this" {
  rule           = aws_cloudwatch_event_rule.this.name
  event_bus_name = var.event_bus_name
  target_id      = "${var.name}-lambda"
  arn            = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "eventbridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.this.arn
} 
