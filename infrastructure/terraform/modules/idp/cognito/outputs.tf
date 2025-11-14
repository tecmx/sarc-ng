output "user_pool_id" {
  description = "The ID of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.id
}

output "user_pool_arn" {
  description = "The ARN of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.arn
}

output "user_pool_endpoint" {
  description = "The endpoint of the Cognito User Pool"
  value       = aws_cognito_user_pool.main.endpoint
}

output "user_pool_client_id" {
  description = "The ID of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.app_client.id
}

output "user_pool_client_secret" {
  description = "The secret of the Cognito User Pool Client"
  value       = aws_cognito_user_pool_client.app_client.client_secret
  sensitive   = true
}

output "user_pool_domain" {
  description = "The domain of the Cognito User Pool"
  value       = var.domain_name != "" ? aws_cognito_user_pool_domain.main[0].domain : null
}

output "identity_pool_id" {
  description = "The ID of the Cognito Identity Pool"
  value       = var.create_identity_pool ? aws_cognito_identity_pool.main[0].id : null
}

output "user_groups" {
  description = "Map of created user groups"
  value = {
    for k, v in aws_cognito_user_group.groups : k => {
      name        = v.name
      description = v.description
      precedence  = v.precedence
    }
  }
}

output "issuer_url" {
  description = "The issuer URL for JWT validation"
  value       = "https://cognito-idp.${data.aws_region.current.name}.amazonaws.com/${aws_cognito_user_pool.main.id}"
}

output "jwks_uri" {
  description = "The JWKS URI for token validation"
  value       = "https://cognito-idp.${data.aws_region.current.name}.amazonaws.com/${aws_cognito_user_pool.main.id}/.well-known/jwks.json"
}

output "ssm_parameters" {
  description = "SSM parameter names"
  value = {
    user_pool_id  = aws_ssm_parameter.user_pool_id.name
    client_id     = aws_ssm_parameter.user_pool_client_id.name
    client_secret = var.generate_client_secret ? aws_ssm_parameter.user_pool_client_secret[0].name : null
  }
}

