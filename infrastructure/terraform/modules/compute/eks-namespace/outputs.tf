output "namespace" {
  description = "The name of the created Kubernetes namespace"
  value       = kubernetes_namespace.this.metadata[0].name
}

output "namespace_uid" {
  description = "The UID of the created Kubernetes namespace"
  value       = kubernetes_namespace.this.metadata[0].uid
}

output "resource_quota_id" {
  description = "The ID of the created resource quota"
  value       = var.create_resource_quota ? kubernetes_resource_quota.this[0].id : null
}

output "admin_role_id" {
  description = "The ID of the created namespace admin role"
  value       = var.create_namespace_admin_role ? kubernetes_role.namespace_admin[0].id : null
}

output "network_policies" {
  description = "Map of created network policy IDs"
  value = {
    default_deny  = var.enable_network_policies ? kubernetes_network_policy.default_deny[0].id : null
    allow_same_ns = var.enable_network_policies ? kubernetes_network_policy.allow_same_namespace[0].id : null
  }
} 
