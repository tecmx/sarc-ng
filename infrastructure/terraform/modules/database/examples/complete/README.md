# Database Module Examples

This directory contains examples of how to use the database module in the SARC-NG project.

## Complete Example

This example demonstrates:

1. How to create a standard RDS instance (MySQL)
2. How to create an Aurora cluster (MySQL compatible)

## Usage

To run this example:

```bash
# Initialize Terraform
terraform init

# Apply the configuration
terraform apply
```

## Notes

- The example uses the default VPC for demonstration purposes
- In a real environment, you would use a properly configured VPC with private subnets
- Deletion protection is disabled for these examples to make cleanup easier

## Cleanup

When you're done experimenting with these examples, run:

```bash
terraform destroy
``` 