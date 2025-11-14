/**
 * EKS Cluster module - Creates an EKS cluster with managed node groups
 */

# EKS Cluster
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name                         = local.cluster_name
  cluster_version                      = var.cluster_version
  cluster_endpoint_public_access       = var.cluster_endpoint_public_access
  cluster_endpoint_private_access      = var.cluster_endpoint_private_access
  cluster_endpoint_public_access_cidrs = var.cluster_endpoint_public_access_cidrs

  vpc_id     = var.vpc_id
  subnet_ids = var.private_subnets

  # IAM OIDC provider
  enable_irsa = true

  # Node groups
  eks_managed_node_group_defaults = var.node_group_defaults
  eks_managed_node_groups         = local.node_groups

  # Authentication
  manage_aws_auth_configmap = true
  aws_auth_roles            = var.aws_auth_roles
  aws_auth_users            = var.aws_auth_users

  tags = local.tags
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_name
}
