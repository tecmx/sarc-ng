# AWS Infrastructure Modules

Detailed documentation of the AWS infrastructure modules used in SARC-NG.

## Table of Contents

- [Network Module](#network-module)
- [Database Module](#database-module)
- [Compute Modules](#compute-modules)
- [DNS Modules](#dns-modules)
- [Observability Module](#observability-module)

## Network Module

The Network module creates a VPC with public, private, and database subnets for the SARC-NG project.

### Features

- Creates a VPC with public, private, and database subnets
- Configurable CIDR blocks and availability zones
- NAT gateways for private subnet internet access
- Internet Gateway for public subnet access
- Option for single or multiple NAT gateways (environment-dependent)
- Kubernetes-ready subnet tagging for EKS
- VPN gateway option
- Database subnet group creation

### Usage

```hcl
module "network" {
  source = "../modules/network"

  project_name = "sarc"
  environment  = "dev"

  # VPC configuration
  vpc_cidr = "10.0.0.0/16"

  # Availability Zones
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  # Subnet configuration
  private_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnet_cidrs   = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  database_subnet_cidrs = ["10.0.201.0/24", "10.0.202.0/24", "10.0.203.0/24"]

  # Optional VPN gateway
  enable_vpn_gateway = false
}
```

### Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| project_name | Project name used for resource naming | string | n/a | yes |
| environment | Environment name (dev, qa, staging, prod) | string | n/a | yes |
| vpc_cidr | CIDR block for the VPC | string | "10.0.0.0/16" | no |
| availability_zones | List of availability zones to use | list(string) | n/a | yes |
| private_subnet_cidrs | CIDR blocks for private subnets | list(string) | n/a | yes |
| public_subnet_cidrs | CIDR blocks for public subnets | list(string) | n/a | yes |
| database_subnet_cidrs | CIDR blocks for database subnets | list(string) | [] | no |
| enable_vpn_gateway | Whether to enable a VPN gateway | bool | false | no |
| public_subnet_tags | Additional tags for public subnets | map(string) | {} | no |
| private_subnet_tags | Additional tags for private subnets | map(string) | {} | no |
| database_subnet_tags | Additional tags for database subnets | map(string) | {} | no |
| additional_tags | Additional tags to add to all resources | map(string) | {} | no |

### Outputs

| Name | Description |
|------|-------------|
| vpc_id | The ID of the VPC |
| vpc_cidr | The CIDR block of the VPC |
| private_subnets | List of IDs of private subnets |
| private_subnet_cidrs | List of CIDR blocks of private subnets |
| public_subnets | List of IDs of public subnets |
| public_subnet_cidrs | List of CIDR blocks of public subnets |
| database_subnets | List of IDs of database subnets |
| database_subnet_cidrs | List of CIDR blocks of database subnets |
| database_subnet_group_name | Name of database subnet group |
| nat_public_ips | List of public Elastic IPs created for NAT gateways |
| availability_zones | List of availability zones used |

### Notes

- In development environments, a single NAT gateway is used for cost savings
- In production environments, one NAT gateway per AZ is created for high availability
- Subnets are tagged for Kubernetes if you plan to use EKS

## Database Module

This module creates a database infrastructure for the SARC-NG project, supporting both standard RDS instances and Aurora clusters.

### Features

- Supports both standard RDS instances and Aurora clusters
- Configurable database engine and version
- Automatic creation of security groups with configurable access rules
- Database credentials stored in AWS Secrets Manager
- Parameter storage in AWS Systems Manager Parameter Store
- Optional multi-AZ deployment
- Configurable backup retention period
- Encryption enabled by default

### Usage

#### Standard RDS Instance

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

#### Aurora Cluster

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

### Inputs

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

### Outputs

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

## Compute Modules

The Compute modules provide resources for deploying applications on AWS:

### Available Modules

- **ECS Cluster**: Creates an Amazon ECS cluster with capacity providers
- **ECS ALB Service**: Deploys an ECS service with an Application Load Balancer
- **ECS NLB Service**: Deploys an ECS service with a Network Load Balancer
- **EKS Cluster**: Creates an Amazon EKS Kubernetes cluster
- **EKS Namespace**: Manages Kubernetes namespaces in an EKS cluster
- **EKS Helm Release**: Deploys applications via Helm charts in an EKS cluster
- **Lambda HTTP API**: Deploys an AWS Lambda function with API Gateway integration
- **Lambda Event Bridge**: Creates a Lambda function triggered by EventBridge events
- **Lambda SQS Consumer**: Deploys a Lambda function to process SQS messages

## DNS Modules

The DNS modules manage domain name configuration:

### Available Modules

- **Zone**: Creates and manages Route53 hosted zones
- **Record**: Manages DNS records within hosted zones

## Observability Module

The Observability module configures monitoring, logging and alerting:

### Features

- CloudWatch log groups with configurable retention
- Custom metric alarms and thresholds
- CloudWatch dashboards for service monitoring
- Integration with existing log groups
