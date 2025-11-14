variable "user_pool_name" {
  description = "Name of the Cognito User Pool"
  type        = string
}

variable "application_name" {
  description = "Name of the application"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "domain_name" {
  description = "Domain name for the Cognito hosted UI"
  type        = string
  default     = ""
}

variable "mfa_configuration" {
  description = "MFA configuration (OFF, ON, OPTIONAL)"
  type        = string
  default     = "OPTIONAL"
  validation {
    condition     = contains(["OFF", "ON", "OPTIONAL"], var.mfa_configuration)
    error_message = "MFA configuration must be OFF, ON, or OPTIONAL"
  }
}

variable "advanced_security_mode" {
  description = "Advanced security mode (OFF, AUDIT, ENFORCED)"
  type        = string
  default     = "AUDIT"
  validation {
    condition     = contains(["OFF", "AUDIT", "ENFORCED"], var.advanced_security_mode)
    error_message = "Advanced security mode must be OFF, AUDIT, or ENFORCED"
  }
}

variable "allow_admin_create_user_only" {
  description = "Allow only admin to create users"
  type        = bool
  default     = false
}

variable "prevent_destroy" {
  description = "Prevent destruction of the user pool"
  type        = bool
  default     = true
}

variable "generate_client_secret" {
  description = "Generate client secret for the app client"
  type        = bool
  default     = false
}

variable "oauth_flows" {
  description = "OAuth flows to enable"
  type        = list(string)
  default     = ["code", "implicit"]
}

variable "oauth_scopes" {
  description = "OAuth scopes to enable"
  type        = list(string)
  default     = ["email", "openid", "profile", "aws.cognito.signin.user.admin"]
}

variable "callback_urls" {
  description = "List of allowed callback URLs"
  type        = list(string)
  default     = []
}

variable "logout_urls" {
  description = "List of allowed logout URLs"
  type        = list(string)
  default     = []
}

variable "identity_providers" {
  description = "List of supported identity providers"
  type        = list(string)
  default     = ["COGNITO"]
}

variable "access_token_validity" {
  description = "Access token validity in minutes"
  type        = number
  default     = 60
}

variable "id_token_validity" {
  description = "ID token validity in minutes"
  type        = number
  default     = 60
}

variable "refresh_token_validity" {
  description = "Refresh token validity in days"
  type        = number
  default     = 30
}

variable "read_attributes" {
  description = "List of user pool attributes the app client can read"
  type        = list(string)
  default     = ["email", "email_verified"]
}

variable "write_attributes" {
  description = "List of user pool attributes the app client can write"
  type        = list(string)
  default     = ["email"]
}

variable "explicit_auth_flows" {
  description = "List of authentication flows"
  type        = list(string)
  default = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]
}

variable "user_groups" {
  description = "Map of user groups to create"
  type = map(object({
    description = string
    precedence  = number
    role_arn    = optional(string)
  }))
  default = {}
}

variable "create_identity_pool" {
  description = "Create a Cognito Identity Pool"
  type        = bool
  default     = false
}

variable "allow_unauthenticated_identities" {
  description = "Allow unauthenticated identities in identity pool"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

