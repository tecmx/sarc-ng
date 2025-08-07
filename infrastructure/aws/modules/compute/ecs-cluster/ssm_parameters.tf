/**
 * ECS Cluster module - SSM parameters
 */

# Store cluster information in SSM
resource "aws_ssm_parameter" "cluster_arn" {
  name        = "/${var.project_name}/${var.environment}/compute/ecs_cluster_arn"
  description = "ECS cluster ARN"
  type        = "String"
  value       = module.ecs.cluster_arn
  tags        = local.tags
}

resource "aws_ssm_parameter" "cluster_name" {
  name        = "/${var.project_name}/${var.environment}/compute/ecs_cluster_name"
  description = "ECS cluster name"
  type        = "String"
  value       = module.ecs.cluster_name
  tags        = local.tags
} 
