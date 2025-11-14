/**
 * Database Schema module - Output values
 */

output "schema_name" {
  description = "Name of the created schema"
  value       = mysql_database.schema.name
}

output "user_name" {
  description = "Name of the created user"
  value       = mysql_user.user.user
}

output "db_secret_arn" {
  description = "ARN of the secret containing user credentials"
  value       = aws_secretsmanager_secret.db_user.arn
} 
