/**
 * Example usage of the eks-helm-release module
 */

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

# Deploy nginx-ingress controller via Helm
module "nginx_ingress" {
  source = "../.."

  project_name = "sarc-ng"
  environment  = "dev"

  release_name = "nginx-ingress"
  repository   = "https://kubernetes.github.io/ingress-nginx"
  chart        = "ingress-nginx"
  version      = "4.8.3"

  namespace        = "ingress-nginx"
  create_namespace = true

  values = {
    controller = {
      service = {
        type = "LoadBalancer"
      }
      metrics = {
        enabled = true
      }
    }
  }

  additional_tags = {
    Component = "ingress"
    Example   = "complete"
  }
}

