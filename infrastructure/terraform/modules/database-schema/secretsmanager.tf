/**
 * Database Schema module - Secrets Manager resources
 */

# Store user credentials in Secrets Manager
resource "aws_secretsmanager_secret" "db_user" {
  name        = "${local.name}-${local.schema_name}-credentials"
  description = "Database credentials for ${local.schema_name}"

  tags = local.tags
}

resource "aws_secretsmanager_secret_version" "db_user" {
  secret_id = aws_secretsmanager_secret.db_user.id

  secret_string = jsonencode({
    username = mysql_user.user.user
    password = random_password.user.result
    engine   = var.engine
    host     = var.host
    port     = var.port
    dbname   = mysql_database.schema.name
  })
} 
