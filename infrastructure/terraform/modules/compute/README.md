# Compute Modules

This directory contains Terraform modules for managing various compute resources in AWS.

## Modules

### ECS Modules

- **ecs-cluster**: Creates an ECS cluster for running containerized applications
- **ecs-alb-service**: Deploys an ECS service with an Application Load Balancer
- **ecs-nlb-service**: Deploys an ECS service with a Network Load Balancer

### EKS Modules

- **eks-cluster**: Creates an Amazon EKS cluster for running Kubernetes workloads
- **eks-namespace**: Creates a Kubernetes namespace with optional resource quotas and RBAC
- **eks-helm-release**: Deploys a Helm chart to a Kubernetes cluster

### Lambda Modules

- **lambda-http-api**: Creates a Lambda function with API Gateway integration
- **lambda-event-bridge**: Creates a Lambda function triggered by EventBridge events
- **lambda-sqs-consumer**: Creates a Lambda function that consumes messages from an SQS queue

## Usage Examples

### ECS Cluster

```hcl
module "ecs_cluster" {
  source = "../../modules/compute/ecs-cluster"

  name        = "app-cluster"
  environment = "dev"
}
```

### ECS Service with ALB

```hcl
module "ecs_service" {
  source = "../../modules/compute/ecs-alb-service"

  name           = "api-service"
  environment    = "dev"
  cluster_id     = module.ecs_cluster.cluster_id
  vpc_id         = module.network.vpc_id
  subnet_ids     = module.network.private_subnet_ids
  container_port = 8080
  
  container_definitions = jsonencode([
    {
      name      = "api"
      image     = "api:latest"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
    }
  ])
}
```

### EKS Cluster

```hcl
module "eks_cluster" {
  source = "../../modules/compute/eks-cluster"

  name                = "app-cluster"
  environment         = "dev"
  vpc_id              = module.network.vpc_id
  subnet_ids          = module.network.private_subnet_ids
  kubernetes_version  = "1.27"
}
```

### Lambda Function with API Gateway

```hcl
module "lambda_api" {
  source = "../../modules/compute/lambda-http-api"

  name        = "api-function"
  handler     = "index.handler"
  runtime     = "nodejs18.x"
  source_path = "../src/api-function.zip"
  
  environment_variables = {
    STAGE = "dev"
  }
}
```

## Module Documentation

For detailed documentation on each module, refer to the README.md file in the respective module directory. 