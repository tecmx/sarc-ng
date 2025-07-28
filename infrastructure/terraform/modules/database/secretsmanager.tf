/**
 * Database secrets management
 */

# Create a secret for database credentials
resource "aws_secretsmanager_secret" "db" {
  name        = "${local.name}-db-credentials"
  description = "Database credentials for ${local.name}"

  tags = local.tags

  recovery_window_in_days = 7
}

# Populate the secret with database credentials
resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id

  secret_string = jsonencode({
    username = var.is_aurora ? module.aurora[0].cluster_master_username : module.db[0].db_instance_username
    password = var.is_aurora ? (
      var.master_password != null ? var.master_password : module.aurora[0].cluster_master_password
      ) : (
      var.master_password != null ? var.master_password : module.db[0].db_instance_password
    )
    engine = var.engine
    host   = var.is_aurora ? module.aurora[0].cluster_endpoint : module.db[0].db_instance_address
    port   = var.is_aurora ? module.aurora[0].cluster_port : module.db[0].db_instance_port
    dbname = local.db_name
  })
} 
