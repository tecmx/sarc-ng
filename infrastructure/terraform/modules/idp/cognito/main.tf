terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Cognito User Pool
resource "aws_cognito_user_pool" "main" {
  name = var.user_pool_name

  # Password policy
  password_policy {
    minimum_length                   = 8
    require_lowercase                = true
    require_uppercase                = true
    require_numbers                  = true
    require_symbols                  = true
    temporary_password_validity_days = 7
  }

  # Account recovery
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  # Auto-verified attributes
  auto_verified_attributes = ["email"]

  # User attributes
  schema {
    attribute_data_type      = "String"
    name                     = "email"
    required                 = true
    mutable                  = true
    developer_only_attribute = false

    string_attribute_constraints {
      min_length = 5
      max_length = 256
    }
  }

  # MFA configuration
  mfa_configuration = var.mfa_configuration

  software_token_mfa_configuration {
    enabled = var.mfa_configuration != "OFF"
  }

  # Admin create user configuration
  admin_create_user_config {
    allow_admin_create_user_only = var.allow_admin_create_user_only

    invite_message_template {
      email_message = "Your username is {username} and temporary password is {####}."
      email_subject = "Your temporary password for ${var.application_name}"
      sms_message   = "Your username is {username} and temporary password is {####}."
    }
  }

  # User pool add-ons
  user_pool_add_ons {
    advanced_security_mode = var.advanced_security_mode
  }

  # Prevent destruction of user pool
  lifecycle {
    prevent_destroy = var.prevent_destroy
  }

  tags = merge(
    var.tags,
    {
      Name        = var.user_pool_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  )
}

# User Pool Domain
resource "aws_cognito_user_pool_domain" "main" {
  count        = var.domain_name != "" ? 1 : 0
  domain       = var.domain_name
  user_pool_id = aws_cognito_user_pool.main.id
}

# User Pool Client for Application
resource "aws_cognito_user_pool_client" "app_client" {
  name         = "${var.application_name}-client"
  user_pool_id = aws_cognito_user_pool.main.id

  generate_secret = var.generate_client_secret

  # OAuth flows
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = var.oauth_flows
  allowed_oauth_scopes                 = var.oauth_scopes

  # Callback URLs
  callback_urls        = var.callback_urls
  logout_urls          = var.logout_urls
  default_redirect_uri = length(var.callback_urls) > 0 ? var.callback_urls[0] : null

  # Supported identity providers
  supported_identity_providers = var.identity_providers

  # Token validity
  access_token_validity  = var.access_token_validity
  id_token_validity      = var.id_token_validity
  refresh_token_validity = var.refresh_token_validity

  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }

  # Read/Write attributes
  read_attributes  = var.read_attributes
  write_attributes = var.write_attributes

  # Prevent user existence errors
  prevent_user_existence_errors = "ENABLED"

  # Enable token revocation
  enable_token_revocation = true

  # Explicit auth flows
  explicit_auth_flows = var.explicit_auth_flows
}

# User Pool Groups
resource "aws_cognito_user_group" "groups" {
  for_each = var.user_groups

  name         = each.key
  user_pool_id = aws_cognito_user_pool.main.id
  description  = each.value.description
  precedence   = each.value.precedence
  role_arn     = try(each.value.role_arn, null)
}

# Identity Pool (optional - for AWS resource access)
resource "aws_cognito_identity_pool" "main" {
  count                            = var.create_identity_pool ? 1 : 0
  identity_pool_name               = "${var.application_name}-identity-pool"
  allow_unauthenticated_identities = var.allow_unauthenticated_identities

  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.app_client.id
    provider_name           = aws_cognito_user_pool.main.endpoint
    server_side_token_check = true
  }

  tags = var.tags
}

# SSM Parameters for application configuration
resource "aws_ssm_parameter" "user_pool_id" {
  name        = "/${var.environment}/${var.application_name}/cognito/user-pool-id"
  description = "Cognito User Pool ID"
  type        = "String"
  value       = aws_cognito_user_pool.main.id

  tags = var.tags
}

resource "aws_ssm_parameter" "user_pool_client_id" {
  name        = "/${var.environment}/${var.application_name}/cognito/client-id"
  description = "Cognito User Pool Client ID"
  type        = "String"
  value       = aws_cognito_user_pool_client.app_client.id

  tags = var.tags
}

resource "aws_ssm_parameter" "user_pool_client_secret" {
  count       = var.generate_client_secret ? 1 : 0
  name        = "/${var.environment}/${var.application_name}/cognito/client-secret"
  description = "Cognito User Pool Client Secret"
  type        = "SecureString"
  value       = aws_cognito_user_pool_client.app_client.client_secret

  tags = var.tags
}

# Data source for current AWS region
data "aws_region" "current" {}

