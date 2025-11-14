# Database Schema Module Example

Creates a database schema with dedicated user and stored credentials.

## Prerequisites

- Existing RDS/Aurora database
- Admin credentials in Secrets Manager
- MySQL provider configured

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- MySQL database schema
- MySQL user with permissions
- Secrets Manager secret with credentials
- SSM parameters for application configuration

## Cleanup

```bash
terraform destroy
```

