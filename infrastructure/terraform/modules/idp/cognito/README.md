# AWS Cognito User Pool Module

Terraform module for creating and managing AWS Cognito User Pool for SARC-NG authentication.

## Features

- Secure password policy (8+ chars, complexity requirements)
- MFA support (TOTP)
- Advanced security mode (compromised credential detection)
- Customizable user groups with precedence
- SSM Parameter Store integration for configuration
- Optional Identity Pool for AWS resource access
- Email verification and account recovery

## Usage

```hcl
module "cognito" {
  source = "../../modules/idp/cognito"

  user_pool_name   = "sarc-ng-dev"
  application_name = "sarc-ng"
  environment      = "dev"
  domain_name      = "sarc-ng-dev"

  # Security settings
  mfa_configuration      = "OPTIONAL"
  advanced_security_mode = "AUDIT"

  # OAuth configuration
  callback_urls = [
    "http://localhost:3000/callback",
    "https://dev.sarc-ng.example.com/callback"
  ]

  logout_urls = [
    "http://localhost:3000",
    "https://dev.sarc-ng.example.com"
  ]

  # User groups
  user_groups = {
    admin = {
      description = "Administrators with full access"
      precedence  = 1
    }
    manager = {
      description = "Managers with resource management access"
      precedence  = 2
    }
    teacher = {
      description = "Teachers who can create classes"
      precedence  = 3
    }
    student = {
      description = "Students who can make reservations"
      precedence  = 4
    }
  }

  tags = {
    Environment = "dev"
    Project     = "sarc-ng"
    ManagedBy   = "Terraform"
  }
}
```

## Outputs

The module provides the following outputs:

- `user_pool_id` - Cognito User Pool ID
- `user_pool_client_id` - App Client ID
- `issuer_url` - JWT issuer URL for validation
- `jwks_uri` - JWKS endpoint for token validation
- `ssm_parameters` - SSM parameter names for configuration

## Configuration in Application

After deploying this module, update your application configuration:

```bash
# Retrieve values from SSM Parameter Store
export COGNITO_USER_POOL_ID=$(aws ssm get-parameter --name "/dev/sarc-ng/cognito/user-pool-id" --query "Parameter.Value" --output text)
export COGNITO_CLIENT_ID=$(aws ssm get-parameter --name "/dev/sarc-ng/cognito/client-id" --query "Parameter.Value" --output text)

# Update aws.env file
echo "export COGNITO_USER_POOL_ID=$COGNITO_USER_POOL_ID" >> aws.env
echo "export COGNITO_CLIENT_ID=$COGNITO_CLIENT_ID" >> aws.env
```

## User Management

### Create a User

```bash
aws cognito-idp admin-create-user \
  --user-pool-id <USER_POOL_ID> \
  --username john.doe \
  --user-attributes Name=email,Value=john.doe@example.com \
  --temporary-password "TempPass123!"
```

### Add User to Group

```bash
aws cognito-idp admin-add-user-to-group \
  --user-pool-id <USER_POOL_ID> \
  --username john.doe \
  --group-name admin
```

### List Users in Group

```bash
aws cognito-idp list-users-in-group \
  --user-pool-id <USER_POOL_ID> \
  --group-name admin
```

## Security Best Practices

1. **Enable MFA**: Set `mfa_configuration = "ON"` for production
2. **Advanced Security**: Use `advanced_security_mode = "ENFORCED"` for production
3. **Token Validity**: Keep access tokens short-lived (default: 60 minutes)
4. **Prevent Destroy**: Set `prevent_destroy = true` to avoid accidental deletion
5. **Callback URLs**: Use HTTPS URLs in production
6. **Client Secret**: Set `generate_client_secret = false` for public clients (SPAs)

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| user_pool_name | Name of the Cognito User Pool | string | n/a | yes |
| application_name | Name of the application | string | n/a | yes |
| environment | Environment name (dev, staging, prod) | string | n/a | yes |
| mfa_configuration | MFA configuration (OFF, ON, OPTIONAL) | string | "OPTIONAL" | no |
| advanced_security_mode | Advanced security mode (OFF, AUDIT, ENFORCED) | string | "AUDIT" | no |
| user_groups | Map of user groups to create | map(object) | {} | no |

See `variables.tf` for complete list of variables.

## Requirements

- Terraform >= 1.0
- AWS Provider ~> 5.0

## License

Proprietary - SARC-NG Project
