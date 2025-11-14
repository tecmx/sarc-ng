/**
 * Example usage of the database-schema module
 */

provider "aws" {
  region = "us-east-1"
}

# Assume database and secrets already exist
data "aws_secretsmanager_secret" "db_admin" {
  name = "/sarc-ng/dev/database/admin"
}

# Create application database schema with dedicated user
module "app_schema" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  admin_secret_arn = data.aws_secretsmanager_secret.db_admin.arn
  endpoint         = "localhost:3306" # Replace with actual RDS endpoint

  schema_name = "sarc_app"
  user_name   = "sarc_app_user"

  additional_tags = {
    Application = "sarc-ng"
    Example     = "complete"
  }
}

