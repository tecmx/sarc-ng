/**
 * EKS Cluster module - SSM parameters
 */

# Store cluster information in SSM
resource "aws_ssm_parameter" "cluster_name" {
  name        = "/${var.project_name}/${var.environment}/compute/eks_cluster_name"
  description = "EKS cluster name"
  type        = "String"
  value       = module.eks.cluster_name
  tags        = local.tags
}

resource "aws_ssm_parameter" "cluster_endpoint" {
  name        = "/${var.project_name}/${var.environment}/compute/eks_cluster_endpoint"
  description = "EKS cluster endpoint"
  type        = "String"
  value       = module.eks.cluster_endpoint
  tags        = local.tags
}

resource "aws_ssm_parameter" "cluster_certificate_authority_data" {
  name        = "/${var.project_name}/${var.environment}/compute/eks_cluster_ca_data"
  description = "EKS cluster certificate authority data"
  type        = "String"
  value       = module.eks.cluster_certificate_authority_data
  tags        = local.tags
}

resource "aws_ssm_parameter" "cluster_oidc_issuer_url" {
  name        = "/${var.project_name}/${var.environment}/compute/eks_cluster_oidc_issuer_url"
  description = "EKS cluster OIDC issuer URL"
  type        = "String"
  value       = module.eks.cluster_oidc_issuer_url
  tags        = local.tags
} 
