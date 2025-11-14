/**
 * Lambda HTTP API module - API Gateway resources
 */

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

