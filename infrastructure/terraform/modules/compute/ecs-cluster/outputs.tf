/**
 * ECS Cluster module - Output values
 */

output "cluster_arn" {
  description = "ARN of the ECS cluster"
  value       = module.ecs.cluster_arn
}

output "cluster_name" {
  description = "Name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "cluster_id" {
  description = "ID of the ECS cluster"
  value       = module.ecs.cluster_id
} 
