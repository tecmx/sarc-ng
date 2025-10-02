/**
 * Database Schema module - Creates a MySQL schema and a user with permissions
 * for a specific service
 */

# Create the database schema and user using a provider
provider "mysql" {
  endpoint = "${var.host}:${var.port}"
  username = jsondecode(data.aws_secretsmanager_secret_version.admin.secret_string)["username"]
  password = jsondecode(data.aws_secretsmanager_secret_version.admin.secret_string)["password"]
}

data "aws_secretsmanager_secret_version" "admin" {
  secret_id = var.admin_secret_arn
}

resource "mysql_database" "schema" {
  name                  = local.schema_name
  default_character_set = "utf8mb4"
  default_collation     = "utf8mb4_unicode_ci"
}

resource "random_password" "user" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "mysql_user" "user" {
  user               = local.user_name
  host               = "%"
  plaintext_password = random_password.user.result
}

resource "mysql_grant" "user" {
  user       = mysql_user.user.user
  host       = mysql_user.user.host
  database   = mysql_database.schema.name
  privileges = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "REFERENCES", "INDEX", "ALTER", "CREATE TEMPORARY TABLES", "LOCK TABLES", "EXECUTE", "CREATE VIEW", "SHOW VIEW", "CREATE ROUTINE", "ALTER ROUTINE"]
}
