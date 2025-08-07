# Infrastructure

Infrastructure management and AWS services overview for SARC-NG.

## Overview

SARC-NG infrastructure is managed using Infrastructure as Code (IaC) principles with Terraform and Terragrunt. The system supports multi-environment deployments with consistent configurations across development, staging, and production environments.

## Architecture Principles

### Account-Based Organization

Infrastructure is organized by AWS accounts for environment isolation:
- **Development Account**: Rapid iteration and testing
- **Staging Account**: Pre-production validation
- **Production Account**: Live system with high availability

### Multi-Region Support

Each account supports multiple AWS regions:
- **Primary Region**: us-east-1 (production workloads)
- **Secondary Region**: us-west-2 (disaster recovery)

### Modular Design

Reusable Terraform modules for consistent deployments:
- **Network**: VPC, subnets, routing, security groups
- **Compute**: ECS clusters, Lambda functions, Auto Scaling Groups
- **Database**: RDS instances, Aurora clusters, backup policies
- **DNS**: Route53 zones, records, SSL certificates
- **Observability**: CloudWatch, monitoring, alerting

## Directory Structure

```
infrastructure/aws/

├── Makefile                   # Removed for simplification
├── live/                      # Environment configurations
│   ├── terragrunt.hcl         # Root configuration
│   └── accounts/              # Account-based organization
│       ├── dev/               # Development environment
│       │   ├── account.hcl    # Account-specific settings
│       │   ├── env.hcl        # Environment variables
│       │   ├── terragrunt.hcl # LocalStack configuration
│       │   └── us-east-1/     # Primary region
│       │       ├── region.hcl # Region-specific settings
│       │       ├── network/   # VPC and networking
│       │       ├── database/  # RDS and Aurora
│       │       ├── compute/   # ECS and Lambda
│       │       ├── dns/       # Route53 configuration
│       │       ├── observability/ # Monitoring setup
│       │       └── services/  # Application services
│       │           ├── sarcng-api/
│       │           │   ├── ecs-alb-service/
│       │           │   ├── lambda-http-api/
│       │           │   ├── dns/record/
│       │           │   └── schema/
│       │           └── sarcng-web/
│       ├── staging/           # Staging environment
│       └── prod/              # Production environment
└── modules/                   # Reusable Terraform modules
    ├── network/               # VPC, subnets, gateways
    ├── compute/               # ECS, EKS, Lambda modules
    │   ├── ecs-cluster/
    │   ├── ecs-alb-service/
    │   ├── ecs-nlb-service/
    │   ├── eks-cluster/
    │   ├── eks-namespace/
    │   ├── eks-helm-release/
    │   ├── lambda-http-api/
    │   ├── lambda-event-bridge/
    │   └── lambda-sqs-consumer/
    ├── database/              # RDS, Aurora modules
    ├── dns/                   # Route53 modules
    │   ├── zone/
    │   └── record/
    └── observability/         # CloudWatch, monitoring
```

## Infrastructure Components

### Network Infrastructure

**VPC Configuration:**
- CIDR blocks sized for environment needs
- Multi-AZ deployment for high availability
- Public, private, and database subnet tiers
- NAT gateways for private subnet internet access

**Development Environment:**
```hcl
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b"]
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24"]
database_subnet_cidrs = ["10.0.20.0/24", "10.0.21.0/24"]
single_nat_gateway = true  # Cost optimization
```

**Production Environment:**
```hcl
vpc_cidr = "10.0.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
public_subnet_cidrs = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24", "10.0.12.0/24"]
database_subnet_cidrs = ["10.0.20.0/24", "10.0.21.0/24", "10.0.22.0/24"]
one_nat_gateway_per_az = true  # High availability
```

### Compute Infrastructure

**ECS Configuration:**
- Fargate for serverless container execution
- Application Load Balancer for traffic distribution
- Auto Scaling based on CPU and memory metrics
- CloudWatch logging for observability

**Lambda Configuration:**
- API Gateway integration for serverless HTTP API
- EventBridge triggers for event-driven processing
- SQS integration for message queue processing
- VPC configuration for database access

**Service Scaling:**
```hcl
# Development
desired_count = 1
cpu = 512
memory = 1024

# Production
desired_count = 3
cpu = 1024
memory = 2048
auto_scaling_enabled = true
```

### Database Infrastructure

**Development Database:**
```hcl
engine = "mysql"
engine_version = "8.0"
instance_class = "db.t3.medium"
allocated_storage = 20
max_allocated_storage = 100
multi_az = false
backup_retention_period = 7
deletion_protection = false
```

**Production Database:**
```hcl
# Aurora MySQL Cluster
is_aurora = true
engine = "aurora-mysql"
engine_version = "8.0.mysql_aurora.3.03.1"
instance_class = "db.r5.large"
instances = {
  1 = { instance_class = "db.r5.large", promotion_tier = 1 }
  2 = { instance_class = "db.r5.large", promotion_tier = 2 }
}
backup_retention_period = 30
deletion_protection = true
performance_insights_enabled = true
```

### DNS and SSL

**Route53 Configuration:**
- Hosted zones for domain management
- Health checks for failover
- Geolocation routing for global users

**SSL Certificate Management:**
- AWS Certificate Manager integration
- Automatic certificate renewal
- Multi-domain and wildcard support

### Security Infrastructure

**Security Groups:**
- Application Load Balancer: 80/443 from internet
- ECS Tasks: Application port from ALB only
- RDS: 3306 from ECS tasks only
- Lambda: VPC access for database connections

**IAM Roles and Policies:**
- ECS Task Execution Role: ECR and CloudWatch access
- ECS Task Role: Application-specific permissions
- Lambda Execution Role: VPC and database access
- Cross-account roles for CI/CD pipelines

## Environment Management

### Local Development with LocalStack

**Prerequisites:**
```bash
# Install LocalStack
pip install localstack

# Start LocalStack services
localstack start -d

# Verify services
curl -s http://localhost:4566/_localstack/health
```

**Infrastructure Testing:**
```bash
# Set LocalStack mode
export LOCALSTACK=true

# Test network module
cd live/accounts/dev/us-east-1/network
terragrunt apply --auto-approve

# Test database module
cd ../database
terragrunt apply --auto-approve

# Clean up
terragrunt destroy --auto-approve
```

### Development Environment

**Network Foundation:**
```bash
cd infrastructure/aws/live/accounts/dev/us-east-1

# Deploy VPC and networking
cd network
terragrunt init
terragrunt plan
terragrunt apply
```

**Database Setup:**
```bash
cd ../database
terragrunt init
terragrunt plan
terragrunt apply
```

**Compute Resources:**
```bash
cd ../compute/ecs-cluster
terragrunt init
terragrunt plan
terragrunt apply
```

**Application Services:**
```bash
cd ../../services/sarcng-api/ecs-alb-service
terragrunt init
terragrunt plan
terragrunt apply
```

### Production Environment

**Bulk Deployment:**
```bash
cd infrastructure/aws/live/accounts/prod/us-east-1

# Plan entire environment
terragrunt run-all plan

# Apply with dependencies
terragrunt run-all apply
```

**Selective Deployment:**
```bash
# Update specific service
cd services/sarcng-api/ecs-alb-service
terragrunt apply

# Update database configuration
cd ../../../database
terragrunt apply
```

## Deployment Patterns

### Blue-Green Deployment

**ECS Blue-Green Strategy:**
1. Deploy new task definition (Green)
2. Update target group with new tasks
3. Monitor health and metrics
4. Switch ALB traffic to Green
5. Terminate Blue tasks after validation

**Lambda Blue-Green Strategy:**
1. Deploy new Lambda version
2. Create alias pointing to new version
3. Gradually shift traffic using weighted routing
4. Monitor error rates and latency
5. Complete cutover or rollback

### Rolling Deployment

**ECS Rolling Updates:**
```hcl
deployment_configuration {
  maximum_percent = 200
  minimum_healthy_percent = 50
}
```

**Auto Scaling Integration:**
```hcl
auto_scaling_target_value = 70.0
scale_in_cooldown = 300
scale_out_cooldown = 300
```

## Monitoring and Observability

### CloudWatch Integration

**Log Groups:**
```hcl
log_groups = [
  "/ecs/sarc-ng/api",
  "/aws/lambda/sarc-ng-api",
  "/aws/rds/instance/sarc-ng-db/error",
  "/aws/rds/instance/sarc-ng-db/slowquery"
]
```

**Custom Metrics:**
```hcl
custom_metrics = [
  "reservation_success_rate",
  "resource_utilization",
  "api_response_time",
  "database_connection_pool"
]
```

**Alerting Rules:**
```hcl
alert_conditions = [
  {
    metric = "CPUUtilization"
    threshold = 80
    comparison = "GreaterThanThreshold"
  },
  {
    metric = "DatabaseConnections"
    threshold = 80
    comparison = "GreaterThanThreshold"
  },
  {
    metric = "HTTPCode_ELB_5XX_Count"
    threshold = 10
    comparison = "GreaterThanThreshold"
  }
]
```

### Performance Monitoring

**Application Performance:**
- Request latency tracking
- Error rate monitoring
- Throughput measurement
- Resource utilization analysis

**Database Performance:**
- Query performance insights
- Connection pool monitoring
- Slow query identification
- Lock contention analysis

## Security Configuration

### Network Security

**VPC Security:**
- Private subnets for application servers
- Database subnets with no internet access
- Network ACLs for additional security layer
- VPC Flow Logs for traffic analysis

**Security Group Rules:**
```hcl
# Application Load Balancer
ingress {
  from_port = 80
  to_port = 80
  protocol = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}

ingress {
  from_port = 443
  to_port = 443
  protocol = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}

# ECS Tasks
ingress {
  from_port = 8080
  to_port = 8080
  protocol = "tcp"
  security_groups = [aws_security_group.alb.id]
}

# RDS Database
ingress {
  from_port = 3306
  to_port = 3306
  protocol = "tcp"
  security_groups = [aws_security_group.ecs_tasks.id]
}
```

### Data Security

**Encryption:**
- RDS encryption at rest using KMS
- EBS volume encryption for ECS instances
- S3 bucket encryption for backups
- SSL/TLS for data in transit

**Secrets Management:**
- AWS Secrets Manager for database credentials
- Systems Manager Parameter Store for configuration
- IAM roles for service authentication
- Regular credential rotation

## Backup and Disaster Recovery

### Database Backup Strategy

**Automated Backups:**
```hcl
backup_retention_period = 30
backup_window = "03:00-06:00"
maintenance_window = "Mon:00:00-Mon:03:00"
copy_tags_to_snapshot = true
```

**Cross-Region Replication:**
```hcl
# Production Aurora cluster
global_cluster_identifier = "sarc-ng-global"
source_region = "us-east-1"
backup_region = "us-west-2"
```

### Application Backup

**Container Registry:**
- Multiple registry copies in different regions
- Tagged versions for rollback capability
- Automated vulnerability scanning

**Configuration Backup:**
- Terraform state in versioned S3 buckets
- Cross-region state replication
- Configuration drift detection

## Cost Optimization

### Development Environment

**Cost-Saving Measures:**
- Single NAT gateway instead of one per AZ
- Smaller instance sizes for non-critical workloads
- Spot instances for development ECS tasks
- Scheduled scaling to reduce off-hours costs

### Production Environment

**Efficiency Strategies:**
- Reserved instances for predictable workloads
- Auto Scaling for dynamic capacity management
- CloudWatch cost monitoring and alerts
- Resource tagging for cost allocation

### Cost Monitoring

**Resource Tagging:**
```hcl
default_tags = {
  Project = "sarc-ng"
  Environment = var.environment
  CostCenter = "engineering"
  Owner = "platform-team"
  ManagedBy = "terraform"
}
```

## Troubleshooting

### Common Issues

**Terraform State Conflicts:**
```bash
# View current state
terragrunt state list

# Remove conflicting resource
terragrunt state rm aws_instance.example

# Import existing resource
terragrunt import aws_instance.example i-1234567890abcdef0
```

**LocalStack Connection Issues:**
```bash
# Check LocalStack health
curl -s http://localhost:4566/_localstack/health

# Restart LocalStack
localstack stop
localstack start -d

# Clear LocalStack data
localstack stop
rm -rf ~/.cache/localstack
localstack start -d
```

**VPC Resource Conflicts:**
```bash
# Check existing VPCs
aws ec2 describe-vpcs --region us-east-1

# Check subnet conflicts
aws ec2 describe-subnets --region us-east-1

# Verify security group rules
aws ec2 describe-security-groups --region us-east-1
```

### Performance Issues

**ECS Service Scaling:**
```bash
# Check service status
aws ecs describe-services \
  --cluster sarc-ng-cluster \
  --services sarc-ng-api

# Update desired count
aws ecs update-service \
  --cluster sarc-ng-cluster \
  --service sarc-ng-api \
  --desired-count 5
```

**Database Performance:**
```bash
# Check RDS metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/RDS \
  --metric-name CPUUtilization \
  --dimensions Name=DBInstanceIdentifier,Value=sarc-ng-db

# Enable Performance Insights
aws rds modify-db-instance \
  --db-instance-identifier sarc-ng-db \
  --enable-performance-insights
```

## Best Practices

### Infrastructure Development

**Module Design:**
- Keep modules focused and reusable
- Use consistent naming conventions
- Document input and output variables
- Include examples in module README

**State Management:**
- Use remote state with S3 and DynamoDB
- Enable state locking to prevent conflicts
- Regular state backup to different regions
- Environment-specific state buckets

**Security Practices:**
- Never commit sensitive data to version control
- Use IAM roles instead of access keys
- Enable CloudTrail for audit logging
- Regular security assessments and updates

### Deployment Workflow

**Change Management:**
1. Test in LocalStack for rapid iteration
2. Deploy to development environment
3. Run integration tests and validation
4. Promote to staging for final testing
5. Deploy to production with monitoring

**Version Control:**
- Tag infrastructure releases
- Use semantic versioning for modules
- Maintain changelog for major updates
- Branch protection for production code
