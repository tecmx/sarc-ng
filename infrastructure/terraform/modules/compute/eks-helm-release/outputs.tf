output "metadata" {
  description = "Block status of the deployed release"
  value       = helm_release.this.metadata
}

output "name" {
  description = "Name of the release"
  value       = helm_release.this.name
}

output "namespace" {
  description = "Namespace of the release"
  value       = helm_release.this.namespace
}

output "version" {
  description = "Version of the release"
  value       = helm_release.this.version
}

output "status" {
  description = "Status of the release"
  value       = helm_release.this.status
}

output "chart" {
  description = "The chart name that was deployed"
  value       = helm_release.this.chart
}

output "chart_version" {
  description = "The chart version that was deployed"
  value       = helm_release.this.version
}

output "app_version" {
  description = "The app version that was deployed"
  value       = try(helm_release.this.metadata[0].app_version, null)
} 
