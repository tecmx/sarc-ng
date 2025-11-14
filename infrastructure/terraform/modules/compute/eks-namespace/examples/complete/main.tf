/**
 * Example usage of the eks-namespace module
 */

provider "kubernetes" {
  config_path = "~/.kube/config"
}

# Create Kubernetes namespace with resource quotas
module "app_namespace" {
  source = "../.."

  namespace = "sarc-ng-app"
  environment = "dev"

  # Resource quotas
  resource_quotas = {
    "requests.cpu"    = "4"
    "requests.memory" = "8Gi"
    "limits.cpu"      = "8"
    "limits.memory"   = "16Gi"
    "pods"            = "20"
  }

  # Network policy (deny all by default)
  create_network_policy = true

  labels = {
    app     = "sarc-ng"
    team    = "engineering"
    example = "complete"
  }
}

