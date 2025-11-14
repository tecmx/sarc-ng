# ECS NLB Service Module Example

Deploys an ECS Fargate service with Network Load Balancer (ideal for gRPC or TCP services).

## Prerequisites

- Existing VPC with subnets
- ECS cluster
- Container image

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- ECS Service (Fargate)
- Network Load Balancer
- Target Group
- Security group for ECS tasks
- CloudWatch log group

## Cleanup

```bash
terraform destroy
```

