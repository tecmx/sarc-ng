/**
 * Observability module - Creates Prometheus and Grafana for monitoring
 */

# Kubernetes provider configuration for connecting to the EKS cluster
provider "kubernetes" {
  host                   = var.cluster_endpoint
  cluster_ca_certificate = base64decode(var.cluster_ca_certificate)
  token                  = var.cluster_token

  # Use exec block for authentication when token isn't provided
  dynamic "exec" {
    for_each = var.cluster_token == null ? [true] : []
    content {
      api_version = "client.authentication.k8s.io/v1beta1"
      command     = "aws"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
    }
  }
}

# Helm provider configuration for deploying Prometheus and Grafana charts
provider "helm" {
  kubernetes {
    host                   = var.cluster_endpoint
    cluster_ca_certificate = base64decode(var.cluster_ca_certificate)
    token                  = var.cluster_token

    # Use exec block for authentication when token isn't provided
    dynamic "exec" {
      for_each = var.cluster_token == null ? [true] : []
      content {
        api_version = "client.authentication.k8s.io/v1beta1"
        command     = "aws"
        args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
      }
    }
  }
}

# Create Prometheus namespace
resource "kubernetes_namespace" "prometheus" {
  metadata {
    name = local.prometheus_namespace
    labels = {
      name        = local.prometheus_namespace
      environment = var.environment
      managed-by  = "terraform"
    }
  }
}

# Create Grafana namespace
resource "kubernetes_namespace" "grafana" {
  metadata {
    name = local.grafana_namespace
    labels = {
      name        = local.grafana_namespace
      environment = var.environment
      managed-by  = "terraform"
    }
  }
}

# Install Prometheus using the Prometheus community Helm chart
resource "helm_release" "prometheus" {
  name       = "prometheus"
  repository = var.prometheus_helm_chart_repository
  chart      = "prometheus"
  version    = var.prometheus_helm_chart_version
  namespace  = kubernetes_namespace.prometheus.metadata[0].name

  # Set common configuration values
  values = [
    yamlencode(local.default_prometheus_values)
  ]

  # Set custom scrape configurations if provided
  dynamic "set" {
    for_each = var.prometheus_scrape_configs
    content {
      name  = "server.${set.key}"
      value = set.value
    }
  }

  # Set alert rules if provided
  dynamic "set" {
    for_each = var.alert_rules
    content {
      name  = "serverFiles.alerts.${set.key}.yaml"
      value = set.value
    }
  }
}

# Install Grafana using the Grafana Helm chart
resource "helm_release" "grafana" {
  name       = "grafana"
  repository = var.grafana_helm_chart_repository
  chart      = "grafana"
  version    = var.grafana_helm_chart_version
  namespace  = kubernetes_namespace.grafana.metadata[0].name

  # Set common configuration values
  values = [
    yamlencode(local.default_grafana_values)
  ]

  # Install additional plugins if specified
  dynamic "set" {
    for_each = toset(var.grafana_additional_plugins)
    content {
      name  = "plugins[${set.key}]"
      value = set.value
    }
  }

  # Configure ingress if enabled
  dynamic "set" {
    for_each = var.enable_ingress ? ["true"] : []
    content {
      name  = "ingress.enabled"
      value = "true"
    }
  }

  dynamic "set" {
    for_each = var.enable_ingress && var.ingress_host != "" ? ["true"] : []
    content {
      name  = "ingress.hosts[0]"
      value = var.ingress_host
    }
  }

  dynamic "set" {
    for_each = var.enable_ingress && var.ingress_tls_secret_name != "" ? ["true"] : []
    content {
      name  = "ingress.tls[0].secretName"
      value = var.ingress_tls_secret_name
    }
  }

  dynamic "set" {
    for_each = var.enable_ingress && var.ingress_host != "" && var.ingress_tls_secret_name != "" ? ["true"] : []
    content {
      name  = "ingress.tls[0].hosts[0]"
      value = var.ingress_host
    }
  }

  # Add Ingress annotations if provided
  dynamic "set" {
    for_each = var.enable_ingress ? var.ingress_annotations : {}
    content {
      name  = "ingress.annotations.${set.key}"
      value = set.value
    }
  }

  depends_on = [
    helm_release.prometheus
  ]
}
