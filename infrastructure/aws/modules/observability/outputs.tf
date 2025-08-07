/**
 * Observability module - Output values
 */

output "prometheus_namespace" {
  description = "Namespace where Prometheus is deployed"
  value       = kubernetes_namespace.prometheus.metadata[0].name
}

output "grafana_namespace" {
  description = "Namespace where Grafana is deployed"
  value       = kubernetes_namespace.grafana.metadata[0].name
}

output "prometheus_server_endpoint" {
  description = "Internal endpoint for the Prometheus server"
  value       = "http://prometheus-server.${kubernetes_namespace.prometheus.metadata[0].name}.svc.cluster.local"
}

output "alertmanager_endpoint" {
  description = "Internal endpoint for the Alertmanager"
  value       = var.enable_alertmanager ? "http://prometheus-alertmanager.${kubernetes_namespace.prometheus.metadata[0].name}.svc.cluster.local" : null
}

output "grafana_endpoint" {
  description = "Internal endpoint for Grafana"
  value       = "http://grafana.${kubernetes_namespace.grafana.metadata[0].name}.svc.cluster.local"
}

output "grafana_admin_password" {
  description = "Admin password for Grafana"
  value       = var.grafana_admin_password
  sensitive   = true
}

output "grafana_ingress_host" {
  description = "Hostname for the Grafana ingress, if enabled"
  value       = var.enable_ingress ? var.ingress_host : null
} 
