# Cognito IDP Module Example

Creates a Cognito User Pool with user groups and identity pool.

## Prerequisites

- AWS account with Cognito permissions
- Application URLs for callbacks

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Cognito User Pool
- User Pool Domain
- User Pool Client
- User Groups
- Identity Pool (optional)
- SSM Parameters

## Configuration

After deployment, retrieve configuration:

```bash
aws ssm get-parameter --name "/dev/sarc-ng/cognito/user-pool-id"
aws ssm get-parameter --name "/dev/sarc-ng/cognito/client-id"
```

## Cleanup

```bash
terraform destroy
```

**Note**: Set `prevent_destroy = true` in production to avoid accidental deletion.

