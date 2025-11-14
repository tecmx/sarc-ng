# DNS Zone Module Example

Creates a Route53 hosted zone with optional ACM certificate.

## Prerequisites

- Domain registered (or ability to update NS records)
- AWS account with Route53 permissions

## Usage

```bash
terraform init
terraform plan
terraform apply

# Update domain NS records with the nameservers from output
```

## Resources Created

- Route53 Hosted Zone
- ACM Certificate (if enabled)
- DNS validation records for ACM
- SSM parameters

## Cleanup

```bash
terraform destroy
```

