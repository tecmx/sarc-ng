# Network Module for SARC-NG

This Terraform module creates a VPC with public, private, and database subnets for the SARC-NG project.

## Features

- Creates a VPC with public, private, and database subnets
- Configurable CIDR blocks and availability zones
- NAT gateways for private subnet internet access
- Internet Gateway for public subnet access
- Option for single or multiple NAT gateways (environment-dependent)
- Kubernetes-ready subnet tagging for EKS
- VPN gateway option
- Database subnet group creation

## Usage

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

## Module Structure

The module is organized into the following files:

- `main.tf` - VPC and networking resources
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `locals.tf` - Local variables
- `versions.tf` - Terraform and provider versions

## Inputs

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

## Outputs

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

## Examples

See the [examples](./examples) directory for working examples.

## Notes

- In development environments, a single NAT gateway is used for cost savings
- In production environments, one NAT gateway per AZ is created for high availability
- Subnets are tagged for Kubernetes if you plan to use EKS 