/**
 * Observability module - Standard Grafana dashboards
 */

# Define standard dashboard configuration to be used in Grafana values
locals {
  # Standard dashboard configurations
  standard_dashboards = {
    # Standard dashboards provided as ConfigMaps
    "kubernetes-cluster-monitoring" = {
      gnetId     = 10000 # Kubernetes Cluster Monitoring
      revision   = 1
      datasource = "Prometheus"
    }
    "node-exporter" = {
      gnetId     = 1860 # Node Exporter Full
      revision   = 21
      datasource = "Prometheus"
    }
    "kubernetes-apiserver" = {
      gnetId     = 12006 # Kubernetes API Server
      revision   = 1
      datasource = "Prometheus"
    }
    "prometheus-2-0-overview" = {
      gnetId     = 3662 # Prometheus 2.0 Overview
      revision   = 2
      datasource = "Prometheus"
    }
  }

  # Dashboard providers configuration for Grafana
  dashboard_providers = {
    "dashboardproviders.yaml" = {
      apiVersion = 1
      providers = [
        {
          name            = "default"
          orgId           = 1
          folder          = ""
          type            = "file"
          disableDeletion = false
          editable        = true
          options = {
            path = "/var/lib/grafana/dashboards/default"
          }
        },
        {
          name            = "kubernetes"
          orgId           = 1
          folder          = "Kubernetes"
          type            = "file"
          disableDeletion = false
          editable        = true
          options = {
            path = "/var/lib/grafana/dashboards/kubernetes"
          }
        }
      ]
    }
  }

  # Set this to be included in the grafana_additional_config local variable
  dashboard_config = {
    dashboardProviders = local.dashboard_providers
    dashboards = {
      default = {
        "prometheus-stats" = {
          json = <<-EOT
          {
            "annotations": {
              "list": [
                {
                  "builtIn": 1,
                  "datasource": "-- Grafana --",
                  "enable": true,
                  "hide": true,
                  "iconColor": "rgba(0, 211, 255, 1)",
                  "name": "Annotations & Alerts",
                  "type": "dashboard"
                }
              ]
            },
            "editable": true,
            "gnetId": null,
            "graphTooltip": 0,
            "id": null,
            "links": [],
            "panels": [
              {
                "aliasColors": {},
                "bars": false,
                "dashLength": 10,
                "dashes": false,
                "datasource": "Prometheus",
                "fieldConfig": {
                  "defaults": {
                    "links": []
                  },
                  "overrides": []
                },
                "fill": 1,
                "fillGradient": 0,
                "gridPos": {
                  "h": 8,
                  "w": 12,
                  "x": 0,
                  "y": 0
                },
                "hiddenSeries": false,
                "id": 2,
                "legend": {
                  "alignAsTable": true,
                  "avg": true,
                  "current": true,
                  "max": true,
                  "min": false,
                  "rightSide": false,
                  "show": true,
                  "total": false,
                  "values": true
                },
                "lines": true,
                "linewidth": 1,
                "nullPointMode": "null",
                "options": {
                  "alertThreshold": true
                },
                "percentage": false,
                "pluginVersion": "7.5.5",
                "pointradius": 2,
                "points": false,
                "renderer": "flot",
                "seriesOverrides": [],
                "spaceLength": 10,
                "stack": false,
                "steppedLine": false,
                "targets": [
                  {
                    "expr": "sum(rate(prometheus_http_requests_total[1m]))",
                    "interval": "",
                    "legendFormat": "HTTP Requests",
                    "refId": "A"
                  }
                ],
                "thresholds": [],
                "timeRegions": [],
                "title": "Prometheus HTTP Request Rate",
                "tooltip": {
                  "shared": true,
                  "sort": 0,
                  "value_type": "individual"
                },
                "type": "graph",
                "xaxis": {
                  "buckets": null,
                  "mode": "time",
                  "name": null,
                  "show": true,
                  "values": []
                },
                "yaxes": [
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  },
                  {
                    "format": "short",
                    "label": null,
                    "logBase": 1,
                    "max": null,
                    "min": null,
                    "show": true
                  }
                ],
                "yaxis": {
                  "align": false,
                  "alignLevel": null
                }
              }
            ],
            "refresh": "10s",
            "schemaVersion": 27,
            "style": "dark",
            "tags": [],
            "templating": {
              "list": []
            },
            "time": {
              "from": "now-6h",
              "to": "now"
            },
            "timepicker": {},
            "timezone": "",
            "title": "Prometheus Stats",
            "uid": "prometheus-stats",
            "version": 1
          }
          EOT
        }
      }
    }
  }

  # Merge standard dashboard configuration with Grafana values
  grafana_additional_config = {
    dashboardProviders = local.dashboard_providers
    dashboards = merge(
      local.dashboard_config.dashboards,
      {
        default = local.standard_dashboards
      }
    )
  }
} 
