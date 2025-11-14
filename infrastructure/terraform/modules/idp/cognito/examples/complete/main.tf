/**
 * Example usage of the Cognito IDP module
 */

provider "aws" {
  region = "us-east-1"
}

# Create Cognito User Pool for authentication
module "auth" {
  source = "../.."

  user_pool_name   = "sarc-ng-dev-users"
  application_name = "sarc-ng"
  environment      = "dev"

  # Hosted UI domain
  domain_name = "sarc-ng-dev-auth"

  # Security settings
  mfa_configuration      = "OPTIONAL"
  advanced_security_mode = "AUDIT"
  prevent_destroy        = false # true for production

  # OAuth configuration
  oauth_flows  = ["code", "implicit"]
  oauth_scopes = ["email", "openid", "profile"]

  callback_urls = [
    "http://localhost:3000/callback",
    "https://dev.sarc-ng.example.com/callback"
  ]

  logout_urls = [
    "http://localhost:3000",
    "https://dev.sarc-ng.example.com"
  ]

  # User groups with permissions
  user_groups = {
    admin = {
      description = "System administrators"
      precedence  = 1
    }
    teacher = {
      description = "Teachers who manage classes"
      precedence  = 2
    }
    student = {
      description = "Students who make reservations"
      precedence  = 3
    }
  }

  # Create Identity Pool for AWS access
  create_identity_pool = true

  tags = {
    Example = "complete"
  }
}

