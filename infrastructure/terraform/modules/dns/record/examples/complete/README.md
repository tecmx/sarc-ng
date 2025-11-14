# DNS Record Module Example

Creates Route53 DNS records (A, CNAME, ALIAS, etc.).

## Prerequisites

- Existing Route53 hosted zone
- Target resources (ALB, CloudFront, etc.)

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Route53 DNS Records (A, CNAME, ALIAS)

## Cleanup

```bash
terraform destroy
```

