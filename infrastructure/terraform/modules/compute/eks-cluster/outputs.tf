/**
 * EKS Cluster module - Output values
 */

output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = module.eks.cluster_name
}

output "cluster_endpoint" {
  description = "Endpoint for the EKS cluster API server"
  value       = module.eks.cluster_endpoint
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks.cluster_certificate_authority_data
}

output "cluster_oidc_issuer_url" {
  description = "URL of the OIDC provider for the cluster"
  value       = module.eks.cluster_oidc_issuer_url
}

output "cluster_oidc_provider_arn" {
  description = "ARN of the OIDC provider for the cluster"
  value       = module.eks.oidc_provider_arn
}

output "cluster_security_group_id" {
  description = "ID of the security group created for the EKS cluster control plane"
  value       = module.eks.cluster_security_group_id
}

output "kubeconfig" {
  description = "Kubernetes configuration for connecting to the cluster"
  value = {
    host                   = module.eks.cluster_endpoint
    token                  = data.aws_eks_cluster_auth.cluster.token
    cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
  }
  sensitive = true
} 
