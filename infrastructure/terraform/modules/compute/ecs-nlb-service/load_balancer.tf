/**
 * ECS NLB Service module - Load balancer resources
 */

# Network Load Balancer
resource "aws_lb" "this" {
  name               = substr(local.full_name, 0, 32)
  internal           = false
  load_balancer_type = "network"
  subnets            = var.public_subnets

  enable_deletion_protection = var.environment == "prod"

  tags = local.tags
}

# Target Group
resource "aws_lb_target_group" "this" {
  name        = substr(local.full_name, 0, 32)
  port        = var.container_port
  protocol    = "TCP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    protocol            = "HTTP"
    path                = var.health_check_path
    port                = var.health_check_port > 0 ? var.health_check_port : "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    interval            = 30
  }

  tags = local.tags
}

# TCP Listener
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = var.listener_port
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
} 
