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
# API Gateway v2
################

resource "aws_apigatewayv2_api" "this" {
  count = var.create_api_gateway ? 1 : 0

  name          = local.api_name
  description   = local.api_description
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins     = local.cors_configuration.allow_origins
    allow_methods     = local.cors_configuration.allow_methods
    allow_headers     = local.cors_configuration.allow_headers
    expose_headers    = local.cors_configuration.expose_headers
    allow_credentials = local.cors_configuration.allow_credentials
    max_age           = local.cors_configuration.max_age
  }

  tags = local.tags
}

resource "aws_apigatewayv2_integration" "this" {
  count = var.create_api_gateway ? 1 : 0

  api_id                 = aws_apigatewayv2_api.this[0].id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.this.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "this" {
  count = var.create_api_gateway ? 1 : 0

  api_id    = aws_apigatewayv2_api.this[0].id
  route_key = "ANY /${var.api_path}"
  target    = "integrations/${aws_apigatewayv2_integration.this[0].id}"
}

resource "aws_apigatewayv2_stage" "this" {
  count = var.create_api_gateway ? 1 : 0

  api_id      = aws_apigatewayv2_api.this[0].id
  name        = var.api_stage_name
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway[0].arn
    format = jsonencode({
      requestId          = "$context.requestId"
      ip                 = "$context.identity.sourceIp"
      requestTime        = "$context.requestTime"
      httpMethod         = "$context.httpMethod"
      routeKey           = "$context.routeKey"
      status             = "$context.status"
      protocol           = "$context.protocol"
      responseLength     = "$context.responseLength"
      path               = "$context.path"
      integrationLatency = "$context.integrationLatency"
      responseLatency    = "$context.responseLatency"
    })
  }

  tags = local.tags
}

resource "aws_cloudwatch_log_group" "api_gateway" {
  count = var.create_api_gateway ? 1 : 0

  name              = "/aws/apigateway/${local.api_name}"
  retention_in_days = var.log_retention_in_days
  tags              = local.tags
}

resource "aws_lambda_permission" "api_gateway" {
  count = var.create_api_gateway ? 1 : 0

  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.this[0].execution_arn}/*/*/${var.api_path}"
}

################
# Custom Domain
################

resource "aws_apigatewayv2_domain_name" "this" {
  count = var.create_api_gateway && var.create_custom_domain && var.api_domain_name != null && var.api_certificate_arn != null ? 1 : 0

  domain_name = var.api_domain_name

  domain_name_configuration {
    certificate_arn = var.api_certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }

  tags = local.tags
}

resource "aws_apigatewayv2_api_mapping" "this" {
  count = var.create_api_gateway && var.create_custom_domain && var.api_domain_name != null && var.api_certificate_arn != null ? 1 : 0

  api_id      = aws_apigatewayv2_api.this[0].id
  domain_name = aws_apigatewayv2_domain_name.this[0].domain_name
  stage       = aws_apigatewayv2_stage.this[0].name
} 
