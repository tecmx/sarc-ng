/**
 * Database subnet groups
 */

# DB subnet group
resource "aws_db_subnet_group" "this" {
  name        = "${local.name}-db-subnet-group"
  description = "DB subnet group for ${local.name}"
  subnet_ids  = var.subnet_ids

  tags = merge(
    local.tags,
    {
      Name = "${local.name}-db-subnet-group"
    }
  )
} 
