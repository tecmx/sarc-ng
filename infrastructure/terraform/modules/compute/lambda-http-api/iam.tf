/**
 * Lambda HTTP API module - IAM resources
 */

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

