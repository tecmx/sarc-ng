/**
 * Example usage of the observability module
 */

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

# Deploy Prometheus and Grafana for EKS monitoring
module "monitoring" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  # EKS cluster configuration
  cluster_name     = "sarc-ng-dev-eks"
  cluster_endpoint = "https://ABC123.gr7.us-east-1.eks.amazonaws.com"
  cluster_ca_cert  = base64decode("LS0tLS1CRUd...")

  # Prometheus configuration
  prometheus_enabled    = true
  prometheus_retention  = "15d"
  prometheus_storage_gb = 50

  # Grafana configuration
  grafana_enabled      = true
  grafana_admin_password = "ChangeMe123!"  # Use secrets manager in prod

  # CloudWatch integration
  enable_cloudwatch_integration = true

  additional_tags = {
    Example = "complete"
  }
}

