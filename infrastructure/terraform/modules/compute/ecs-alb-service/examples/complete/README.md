# ECS ALB Service Module Example

Deploys an ECS Fargate service with Application Load Balancer.

## Prerequisites

- Existing VPC with public and private subnets
- ECS cluster
- Container image in ECR or Docker Hub

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- ECS Service (Fargate)
- Application Load Balancer
- Target Group with health checks
- Security groups (ALB and ECS tasks)
- Auto-scaling policies
- CloudWatch log group

## Cleanup

```bash
terraform destroy
```

