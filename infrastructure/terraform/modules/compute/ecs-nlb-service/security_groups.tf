/**
 * ECS NLB Service module - Security groups
 */

# Security group for ECS tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "${local.full_name}-ecs-tasks-sg"
  description = "Security group for ${local.full_name} ECS tasks"
  vpc_id      = var.vpc_id

  ingress {
    description = "Allow inbound traffic from NLB"
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # NLB doesn't have a security group, so we need to allow all traffic
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
} 
