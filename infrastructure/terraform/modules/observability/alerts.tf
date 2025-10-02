/**
 * Observability module - Standard Prometheus alerting rules
 */

# These are some common alerts that can be used as defaults or overridden
locals {
  default_alert_rules = {
    "groups" = [
      {
        name = "kubernetes-system-alerts"
        rules = [
          {
            alert = "KubernetesPodCrashLooping"
            expr  = "rate(kube_pod_container_status_restarts_total[15m]) * 60 * 5 > 0"
            for   = "10m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "Pod {{ $labels.namespace }}/{{ $labels.pod }} is crash looping",
              description = "Pod {{ $labels.namespace }}/{{ $labels.pod }} is crash looping ({{ $value }} restarts in the last 15 minutes)"
            }
          },
          {
            alert = "KubernetesPodNotReady"
            expr  = "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown\"}) > 0"
            for   = "10m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "Pod {{ $labels.namespace }}/{{ $labels.pod }} is not ready",
              description = "Pod {{ $labels.namespace }}/{{ $labels.pod }} has been in a non-ready state for more than 10 minutes."
            }
          }
        ]
      },
      {
        name = "node-alerts"
        rules = [
          {
            alert = "HighNodeCPU"
            expr  = "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100) > 80"
            for   = "5m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "Node {{ $labels.instance }} has high CPU usage",
              description = "Node {{ $labels.instance }} CPU usage is above 80% for more than 5 minutes (current value: {{ $value }}%)"
            }
          },
          {
            alert = "HighNodeMemoryUsage"
            expr  = "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 80"
            for   = "5m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "Node {{ $labels.instance }} has high memory usage",
              description = "Node {{ $labels.instance }} memory usage is above 80% for more than 5 minutes (current value: {{ $value }}%)"
            }
          },
          {
            alert = "LowNodeDiskSpace"
            expr  = "node_filesystem_avail_bytes{mountpoint=\"/\"} / node_filesystem_size_bytes{mountpoint=\"/\"} * 100 < 20"
            for   = "5m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "Node {{ $labels.instance }} is low on disk space",
              description = "Node {{ $labels.instance }} has less than 20% disk space available on the / partition (current value: {{ $value }}%)"
            }
          }
        ]
      },
      {
        name = "application-alerts"
        rules = [
          {
            alert = "HighRequestLatency"
            expr  = "http_request_duration_seconds{quantile=\"0.9\"} > 2"
            for   = "5m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "High request latency on {{ $labels.instance }}",
              description = "The service {{ $labels.service }} on {{ $labels.instance }} has a 90th percentile latency above 2 seconds for more than 5 minutes (current value: {{ $value }}s)"
            }
          },
          {
            alert = "HighErrorRate"
            expr  = "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m])) * 100 > 5"
            for   = "5m"
            labels = {
              severity = "warning"
            }
            annotations = {
              summary     = "High HTTP error rate",
              description = "The service has a HTTP error rate above 5% for more than 5 minutes (current value: {{ $value }}%)"
            }
          }
        ]
      }
    ]
  }
} 
