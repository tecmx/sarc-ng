/**
 * ECS NLB Service module - Security groups
 */

# Security group for ECS tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "${local.full_name}-ecs-tasks-sg"
  description = "Security group for ${local.full_name} ECS tasks"
  vpc_id      = var.vpc_id

  tags = local.tags

  lifecycle {
    create_before_destroy = true
  }
}

# Ingress rule for ECS tasks from allowed CIDR blocks
resource "aws_security_group_rule" "ecs_tasks_ingress" {
  type              = "ingress"
  from_port         = var.container_port
  to_port           = var.container_port
  protocol          = "tcp"
  cidr_blocks       = var.allowed_cidr_blocks
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow inbound traffic from allowed CIDR blocks"
}

# Egress rule for ECS tasks
resource "aws_security_group_rule" "ecs_tasks_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.ecs_tasks.id
  description       = "Allow all outbound traffic"
}
