/**
 * ECS ALB Service module - Security groups
 */

# Security group for ALB
resource "aws_security_group" "alb" {
  name        = "${local.full_name}-alb-sg"
  description = "Security group for ${local.full_name} ALB"
  vpc_id      = var.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

# Security group for ECS tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "${local.full_name}-ecs-tasks-sg"
  description = "Security group for ${local.full_name} ECS tasks"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Access from ALB"
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
} 
