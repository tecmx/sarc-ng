/**
 * Database security groups
 */

# Security group for database access
resource "aws_security_group" "db" {
  name        = "${local.name}-db-sg"
  description = "Security group for ${local.name} database"
  vpc_id      = var.vpc_id

  tags = merge(
    local.tags,
    {
      Name = "${local.name}-db-sg"
    }
  )

  lifecycle {
    create_before_destroy = true
  }
}

# Allow access from specified CIDR blocks
resource "aws_security_group_rule" "db_ingress_cidr" {
  count = length(var.allowed_cidr_blocks) > 0 ? 1 : 0

  type              = "ingress"
  from_port         = var.port
  to_port           = var.port
  protocol          = "tcp"
  cidr_blocks       = var.allowed_cidr_blocks
  security_group_id = aws_security_group.db.id
  description       = "Allow database access from specified CIDR blocks"
}

# Allow access from specified security groups
resource "aws_security_group_rule" "db_ingress_sg" {
  count = length(var.allowed_security_group_ids)

  type                     = "ingress"
  from_port                = var.port
  to_port                  = var.port
  protocol                 = "tcp"
  source_security_group_id = var.allowed_security_group_ids[count.index]
  security_group_id        = aws_security_group.db.id
  description              = "Allow database access from security group ${var.allowed_security_group_ids[count.index]}"
}

# Allow all outbound traffic
resource "aws_security_group_rule" "db_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.db.id
  description       = "Allow all outbound traffic"
}
