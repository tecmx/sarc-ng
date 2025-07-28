locals {
  rule_name        = var.rule_name != null ? var.rule_name : "${var.name}-rule"
  rule_description = var.rule_description != null ? var.rule_description : "EventBridge rule for ${var.name} Lambda function"

  lambda_source = var.source_path != null ? {
    path = var.source_path
    } : var.source_bucket != null && var.source_key != null ? {
    s3_bucket = var.source_bucket
    s3_key    = var.source_key
  } : null

  lambda_role_name = "${var.name}-lambda-role"

  tags = merge(
    var.tags,
    {
      Name      = var.name
      ManagedBy = "terraform"
    }
  )
} 
