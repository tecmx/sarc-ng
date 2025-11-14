/**
 * Database module - Output values
 */

# RDS/Aurora instance or cluster outputs
output "endpoint" {
  description = "Database endpoint"
  value       = var.is_aurora ? module.aurora[0].cluster_endpoint : module.db[0].db_instance_address
}

output "port" {
  description = "Database port"
  value       = var.is_aurora ? module.aurora[0].cluster_port : module.db[0].db_instance_port
}

output "name" {
  description = "Database name"
  value       = local.db_name
}

output "username" {
  description = "Database master username"
  value       = var.is_aurora ? module.aurora[0].cluster_master_username : module.db[0].db_instance_username
  sensitive   = true
}

# Security group outputs
output "security_group_id" {
  description = "Database security group ID"
  value       = aws_security_group.db.id
}

# Secrets manager outputs
output "credentials_secret_arn" {
  description = "ARN of the Secrets Manager secret containing database credentials"
  value       = aws_secretsmanager_secret.db.arn
}

output "credentials_secret_name" {
  description = "Name of the Secrets Manager secret containing database credentials"
  value       = aws_secretsmanager_secret.db.name
}

# SSM Parameter outputs
output "ssm_parameter_endpoint" {
  description = "SSM parameter name for database endpoint"
  value       = aws_ssm_parameter.db_endpoint.name
}

output "ssm_parameter_port" {
  description = "SSM parameter name for database port"
  value       = aws_ssm_parameter.db_port.name
}

output "ssm_parameter_name" {
  description = "SSM parameter name for database name"
  value       = aws_ssm_parameter.db_name.name
}

output "ssm_parameter_secret_arn" {
  description = "SSM parameter name for database credentials secret ARN"
  value       = aws_ssm_parameter.db_secret_arn.name
}
