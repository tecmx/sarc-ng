variable "name" {
  description = "Name of the Helm release"
  type        = string
}

variable "namespace" {
  description = "Kubernetes namespace to deploy the Helm chart in"
  type        = string
}

variable "repository" {
  description = "Helm repository URL"
  type        = string
  default     = null
}

variable "chart" {
  description = "Helm chart name or local path to chart"
  type        = string
}

variable "chart_version" {
  description = "Specific version of the Helm chart to deploy"
  type        = string
  default     = null
}

variable "values" {
  description = "List of values in raw YAML to pass to Helm"
  type        = list(string)
  default     = []
}

variable "values_files" {
  description = "List of paths to values files to pass to Helm"
  type        = list(string)
  default     = []
}

variable "set" {
  description = "Value block with custom STRING values to be merged with the values yaml"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "set_sensitive" {
  description = "Value block with custom sensitive values to be merged with the values yaml that won't be exposed in the plan's diff"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "create_namespace" {
  description = "Create the Kubernetes namespace if it does not exist"
  type        = bool
  default     = false
}

variable "atomic" {
  description = "If true, installation process purges chart on fail"
  type        = bool
  default     = false
}

variable "cleanup_on_fail" {
  description = "Allow deletion of new resources created in this upgrade when upgrade fails"
  type        = bool
  default     = false
}

variable "timeout" {
  description = "Time in seconds to wait for any individual Kubernetes operation"
  type        = number
  default     = 300
}

variable "recreate_pods" {
  description = "Perform pods restart during upgrade/rollback"
  type        = bool
  default     = false
}

variable "max_history" {
  description = "Maximum number of release versions stored per release"
  type        = number
  default     = 10
}

variable "wait" {
  description = "Will wait until all resources are in a ready state before marking the release as successful"
  type        = bool
  default     = true
}

variable "wait_for_jobs" {
  description = "If wait is enabled, will wait until all Jobs have been completed before marking the release as successful"
  type        = bool
  default     = false
}

variable "reuse_values" {
  description = "When upgrading, reuse the last release's values and merge in any overrides"
  type        = bool
  default     = false
}

variable "force_update" {
  description = "Force resource update through delete/recreate if needed"
  type        = bool
  default     = false
}

variable "dependency_update" {
  description = "Runs helm dependency update before installing the chart"
  type        = bool
  default     = false
}

variable "skip_crds" {
  description = "If set, no CRDs will be installed"
  type        = bool
  default     = false
}

variable "render_subchart_notes" {
  description = "If set, render subchart notes along with the parent"
  type        = bool
  default     = true
}

variable "disable_webhooks" {
  description = "Prevent hooks from running"
  type        = bool
  default     = false
}

variable "replace" {
  description = "Re-use the given name, only if that name is a deleted release which remains in the history"
  type        = bool
  default     = false
}

variable "description" {
  description = "Release description attribute"
  type        = string
  default     = null
}

variable "postrender" {
  description = "Configure a command to run after helm renders the manifest which can alter the manifest contents"
  type = list(object({
    binary_path = string
    args        = list(string)
  }))
  default = []
} 
