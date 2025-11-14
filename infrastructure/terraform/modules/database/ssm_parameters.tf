/**
 * Database SSM parameters
 */

# Store database endpoint in SSM Parameter Store
resource "aws_ssm_parameter" "db_endpoint" {
  name        = "/${var.project_name}/${var.environment}/database/endpoint"
  description = "Database endpoint for ${local.name}"
  type        = "String"
  value       = var.is_aurora ? module.aurora[0].cluster_endpoint : module.db[0].db_instance_address

  tags = local.tags
}

# Store database port in SSM Parameter Store
resource "aws_ssm_parameter" "db_port" {
  name        = "/${var.project_name}/${var.environment}/database/port"
  description = "Database port for ${local.name}"
  type        = "String"
  value       = var.is_aurora ? module.aurora[0].cluster_port : module.db[0].db_instance_port

  tags = local.tags
}

# Store database name in SSM Parameter Store
resource "aws_ssm_parameter" "db_name" {
  name        = "/${var.project_name}/${var.environment}/database/name"
  description = "Database name for ${local.name}"
  type        = "String"
  value       = local.db_name

  tags = local.tags
}

# Store database secret ARN in SSM Parameter Store
resource "aws_ssm_parameter" "db_secret_arn" {
  name        = "/${var.project_name}/${var.environment}/database/secret_arn"
  description = "Database credentials secret ARN for ${local.name}"
  type        = "String"
  value       = aws_secretsmanager_secret.db.arn

  tags = local.tags
} 
