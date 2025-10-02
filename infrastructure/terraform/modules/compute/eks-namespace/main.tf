resource "kubernetes_namespace" "this" {
  metadata {
    name = var.namespace
    labels = merge(
      {
        name        = var.namespace
        environment = var.environment
      },
      var.labels
    )
    annotations = var.annotations
  }
}

resource "kubernetes_resource_quota" "this" {
  count = var.create_resource_quota ? 1 : 0

  metadata {
    name      = "${var.namespace}-quota"
    namespace = kubernetes_namespace.this.metadata[0].name
  }

  spec {
    hard = {
      "requests.cpu"    = var.quota_requests_cpu
      "requests.memory" = var.quota_requests_memory
      "limits.cpu"      = var.quota_limits_cpu
      "limits.memory"   = var.quota_limits_memory
      "pods"            = var.quota_pods
    }
  }
}

resource "kubernetes_role" "namespace_admin" {
  count = var.create_namespace_admin_role ? 1 : 0

  metadata {
    name      = "${var.namespace}-admin"
    namespace = kubernetes_namespace.this.metadata[0].name
  }

  rule {
    api_groups = ["", "extensions", "apps", "batch", "networking.k8s.io"]
    resources  = ["*"]
    verbs      = ["*"]
  }

  rule {
    api_groups = ["rbac.authorization.k8s.io"]
    resources  = ["roles", "rolebindings"]
    verbs      = ["*"]
  }
}

resource "kubernetes_role_binding" "namespace_admin" {
  count = var.create_namespace_admin_role && length(var.namespace_admins) > 0 ? 1 : 0

  metadata {
    name      = "${var.namespace}-admin-binding"
    namespace = kubernetes_namespace.this.metadata[0].name
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = kubernetes_role.namespace_admin[0].metadata[0].name
  }

  dynamic "subject" {
    for_each = var.namespace_admins

    content {
      kind      = "User"
      name      = subject.value
      api_group = "rbac.authorization.k8s.io"
    }
  }
}

resource "kubernetes_network_policy" "default_deny" {
  count = var.enable_network_policies ? 1 : 0

  metadata {
    name      = "${var.namespace}-default-deny"
    namespace = kubernetes_namespace.this.metadata[0].name
  }

  spec {
    pod_selector {}
    policy_types = ["Ingress", "Egress"]
  }
}

resource "kubernetes_network_policy" "allow_same_namespace" {
  count = var.enable_network_policies ? 1 : 0

  metadata {
    name      = "${var.namespace}-allow-same-namespace"
    namespace = kubernetes_namespace.this.metadata[0].name
  }

  spec {
    pod_selector {}

    ingress {
      from {
        namespace_selector {
          match_labels = {
            name = var.namespace
          }
        }
      }
    }

    egress {
      to {
        namespace_selector {
          match_labels = {
            name = var.namespace
          }
        }
      }
    }
  }
} 
