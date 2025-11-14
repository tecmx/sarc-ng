/**
 * Database module - Creates an RDS/Aurora cluster
 */

# Choose between Aurora and standard RDS based on is_aurora variable
module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 6.0"
  count   = var.is_aurora ? 0 : 1

  identifier = "${local.name}-db"

  engine                = var.engine
  engine_version        = var.engine_version
  instance_class        = var.instance_class
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage

  db_name  = local.db_name
  username = "admin"
  # Handle password management - either use provided password or let RDS generate one
  password                    = var.master_password
  manage_master_user_password = var.master_password == null ? true : false

  multi_az               = var.multi_az
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.db.id]

  # Performance and monitoring
  monitoring_interval             = 60
  monitoring_role_name            = "${local.name}-rds-monitoring-role"
  create_monitoring_role          = true
  performance_insights_enabled    = true
  performance_insights_kms_key_id = null # Use default AWS managed key

  # Backup and maintenance
  backup_retention_period = var.backup_retention_period
  backup_window           = "03:00-06:00"
  maintenance_window      = "Mon:00:00-Mon:03:00"

  # Snapshot configuration
  skip_final_snapshot = var.environment != "prod"
  # RDS module uses final_snapshot_identifier_prefix
  final_snapshot_identifier_prefix = "${local.name}-db-snapshot"

  # Security
  deletion_protection = var.deletion_protection
  storage_encrypted   = true

  tags = local.tags
}

# Aurora cluster for MySQL or PostgreSQL
module "aurora" {
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "~> 8.0"
  count   = var.is_aurora ? 1 : 0

  name           = "${local.name}-aurora"
  engine         = var.engine
  engine_version = var.engine_version

  # Default to two instances - one writer, one reader
  instance_class = var.instance_class
  instances = {
    1 = {
      instance_class = var.instance_class
      promotion_tier = 1
    }
    2 = {
      instance_class = var.instance_class
      promotion_tier = 2
    }
  }

  vpc_id                 = var.vpc_id
  subnets                = var.subnet_ids
  create_db_subnet_group = true
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.db.id]

  # Database configuration
  database_name               = local.db_name
  master_username             = "admin"
  manage_master_user_password = var.master_password == null ? true : false
  master_password             = var.master_password

  # Performance and monitoring
  monitoring_interval          = 60
  create_monitoring_role       = true
  performance_insights_enabled = true

  # Backup and maintenance
  backup_retention_period      = var.backup_retention_period
  preferred_backup_window      = "03:00-06:00"
  preferred_maintenance_window = "Mon:00:00-Mon:03:00"

  # Snapshot configuration - Aurora module uses skip_final_snapshot directly
  skip_final_snapshot = var.environment != "prod"
  # For Aurora v8.0, we use final_snapshot_identifier
  final_snapshot_identifier = var.environment == "prod" ? "${local.name}-aurora-snapshot" : null

  # Security
  deletion_protection = var.deletion_protection
  storage_encrypted   = true

  # Additional configuration
  apply_immediately          = true
  auto_minor_version_upgrade = true

  tags = local.tags
}
