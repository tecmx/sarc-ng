# Database Module for SARC-NG

This Terraform module creates a database infrastructure for the SARC-NG project, supporting both standard RDS instances and Aurora clusters.

## Features

- Supports both standard RDS instances and Aurora clusters
- Configurable database engine and version
- Automatic creation of security groups with configurable access rules
- Database credentials stored in AWS Secrets Manager
- Parameter storage in AWS Systems Manager Parameter Store
- Optional multi-AZ deployment
- Configurable backup retention period
- Encryption enabled by default

## Usage

### Standard RDS Instance

```hcl
module "database" {
  source = "../modules/database"

  project_name = "sarc"
  environment  = "dev"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # Database configuration
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = "db.t3.medium"

  # Security
  allowed_cidr_blocks = ["10.0.0.0/16"]
}
```

### Aurora Cluster

```hcl
module "database" {
  source = "../modules/database"

  project_name = "sarc"
  environment  = "prod"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # Use Aurora instead of standard RDS
  is_aurora      = true
  engine         = "aurora-mysql"
  engine_version = "8.0.mysql_aurora.3.03.1"
  instance_class = "db.r5.large"

  # Security
  allowed_cidr_blocks = ["10.0.0.0/16"]
}
```

## Module Structure

The module is organized into the following files:

- `main.tf` - Core database resources (RDS or Aurora)
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `locals.tf` - Local variables
- `versions.tf` - Terraform and provider versions
- `subnet_groups.tf` - DB subnet group configuration
- `security_groups.tf` - Security group configuration
- `secrets.tf` - AWS Secrets Manager resources
- `ssm_parameters.tf` - AWS Systems Manager Parameter Store resources

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| project_name | Project name used for resource naming | string | n/a | yes |
| environment | Environment name (dev, qa, staging, prod) | string | n/a | yes |
| vpc_id | ID of the VPC where database will be deployed | string | n/a | yes |
| subnet_ids | List of subnet IDs for the DB subnet group | list(string) | n/a | yes |
| engine | Database engine (mysql, postgres, aurora-mysql, aurora-postgresql) | string | "mysql" | no |
| engine_version | Database engine version | string | "8.0" | no |
| instance_class | Instance type for the RDS instances | string | "db.t3.medium" | no |
| allocated_storage | Allocated storage in GB (not applicable for Aurora) | number | 20 | no |
| max_allocated_storage | Max allocated storage for autoscaling (not applicable for Aurora) | number | 100 | no |
| is_aurora | Whether to create an Aurora cluster instead of a standard RDS instance | bool | false | no |
| multi_az | Whether to deploy a multi-AZ RDS instance (not applicable for Aurora) | bool | false | no |
| backup_retention_period | Days to retain backups | number | 7 | no |
| deletion_protection | Enable deletion protection | bool | true | no |
| master_password | Master password for the database (if not provided, a random one will be generated) | string | null | no |
| port | Database port (3306 for MySQL, 5432 for PostgreSQL) | number | 3306 | no |
| allowed_cidr_blocks | List of CIDR blocks that are allowed to access the database | list(string) | [] | no |
| allowed_security_group_ids | List of security group IDs that are allowed to access the database | list(string) | [] | no |
| additional_tags | Additional tags to add to all resources | map(string) | {} | no |

## Outputs

| Name | Description |
|------|-------------|
| endpoint | Database endpoint |
| port | Database port |
| name | Database name |
| username | Database master username |
| security_group_id | Database security group ID |
| credentials_secret_arn | ARN of the Secrets Manager secret containing database credentials |
| credentials_secret_name | Name of the Secrets Manager secret containing database credentials |
| ssm_parameter_endpoint | SSM parameter name for database endpoint |
| ssm_parameter_port | SSM parameter name for database port |
| ssm_parameter_name | SSM parameter name for database name |
| ssm_parameter_secret_arn | SSM parameter name for database credentials secret ARN |

## Examples

See the [examples](./examples) directory for working examples. 