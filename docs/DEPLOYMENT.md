# Deployment

Production deployment procedures and options for SARC-NG.

## Deployment Options

SARC-NG supports multiple deployment strategies to accommodate different environments and requirements:

1. **Local Development**: Native Go or Docker Compose
2. **Containerized**: Docker-based deployments
3. **Serverless**: AWS Lambda with API Gateway
4. **Traditional**: Binary deployment with external database
5. **Cloud Native**: ECS or EKS with managed services

## Local Development

### Native Go Development

**Prerequisites:**
- Go 1.24+
- MySQL 8.0+ (local or Docker)

**Quick Start:**
```bash
# Install dependencies
make setup

# Start database (if using Docker)
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=sarcng \
  -p 3306:3306 mysql:8.0

# Run application
make run
```

### Docker Compose Development

**Prerequisites:**
- Docker and Docker Compose

**Quick Start:**
```bash
# Start full development environment
docker compose up -d

# View application
open http://localhost:8080

# Stop and cleanup
docker compose down -v --remove-orphans
```

## Container Deployment

### Docker Image Build

**Development Image:**
```bash
docker build --target development -t sarc-ng:dev .
```

**Production Image:**
```bash
docker build --target production -t sarc-ng:prod .
```

### Container Registry

**Push to Registry:**
```bash
# Tag for registry
docker tag sarc-ng:prod your-registry/sarc-ng:latest

# Push to registry
docker push your-registry/sarc-ng:latest
```

### Standalone Container

**Run Production Container:**
```bash
docker run -d \
  --name sarc-ng \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  sarc-ng:prod
```

## Serverless Deployment

### AWS Lambda with API Gateway

**Build Lambda Package:**
```bash
# Build Lambda binary
GOOS=linux GOARCH=amd64 go build -o lambda ./cmd/lambda

# Create deployment package
zip lambda-deployment.zip lambda
```

**Deploy with AWS CLI:**
```bash
# Create function
aws lambda create-function \
  --function-name sarc-ng-api \
  --runtime provided.al2 \
  --role arn:aws:iam::account:role/lambda-execution-role \
  --handler lambda \
  --zip-file fileb://lambda-deployment.zip

# Create API Gateway (see Terraform modules for complete setup)
```

**Environment Variables for Lambda:**
```bash
DB_HOST=your-rds-endpoint
DB_PASSWORD=your-password
DB_NAME=sarcng
GIN_MODE=release
```

## Infrastructure as Code

### Terraform Deployment

**Prerequisites:**
- Terraform 1.0+
- Terragrunt 0.45+
- AWS CLI configured

**Directory Structure:**
```
infrastructure/terraform/live/accounts/
├── dev/                    # Development environment
├── staging/                # Staging environment  
└── prod/                   # Production environment
```

### Development Environment

**Deploy Development Infrastructure:**
```bash
cd infrastructure/terraform/live/accounts/dev/us-east-1

# Deploy network foundation
cd network
terragrunt plan
terragrunt apply

# Deploy database
cd ../database
terragrunt plan
terragrunt apply

# Deploy compute resources
cd ../compute/ecs-cluster
terragrunt plan
terragrunt apply

# Deploy application service
cd ../../services/sarcng-api/ecs-alb-service
terragrunt plan
terragrunt apply
```

### Production Environment

**Deploy Production Infrastructure:**
```bash
cd infrastructure/terraform/live/accounts/prod/us-east-1

# Deploy all modules
terragrunt run-all plan
terragrunt run-all apply
```

### Local Infrastructure Testing

**Using LocalStack:**
```bash
# Start LocalStack
localstack start -d

# Test infrastructure locally
cd infrastructure/terraform/live/accounts/dev/us-east-1/network
LOCALSTACK=true terragrunt apply --auto-approve
```

## AWS ECS Deployment

### ECS with Application Load Balancer

**Service Configuration:**
```yaml
# Terraform configuration
service_name: "sarc-ng-api"
container_image: "your-registry/sarc-ng:latest"
container_port: 8080
cpu: 512
memory: 1024
desired_count: 2

# Environment variables
environment:
  DB_HOST: "{{ rds_endpoint }}"
  DB_NAME: "sarcng"
  GIN_MODE: "release"
  PORT: "8080"
```

**Deployment Process:**
1. Build and push Docker image to ECR
2. Deploy ECS cluster via Terraform
3. Deploy service with ALB via Terraform
4. Configure domain and SSL certificate
5. Set up monitoring and logging

### ECS Service Update

**Update Service with New Image:**
```bash
# Update task definition with new image
aws ecs update-service \
  --cluster sarc-ng-cluster \
  --service sarc-ng-api \
  --force-new-deployment
```

## AWS EKS Deployment

### Kubernetes Deployment

**Kubernetes Manifests:**
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sarc-ng-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sarc-ng-api
  template:
    metadata:
      labels:
        app: sarc-ng-api
    spec:
      containers:
      - name: api
        image: your-registry/sarc-ng:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
```

**Helm Chart Deployment:**
```bash
# Deploy via Terraform Helm provider
cd infrastructure/terraform/live/accounts/prod/us-east-1/services/sarcng-api/eks-helm-release
terragrunt apply
```

## Database Deployment

### RDS MySQL

**Development Database:**
```hcl
# Terraform configuration
engine = "mysql"
engine_version = "8.0"
instance_class = "db.t3.medium"
allocated_storage = 20
multi_az = false
backup_retention_period = 7
```

**Production Database:**
```hcl
# Terraform configuration
engine = "mysql"
engine_version = "8.0"
instance_class = "db.r5.xlarge"
allocated_storage = 100
multi_az = true
backup_retention_period = 30
deletion_protection = true
```

### Aurora MySQL Cluster

**Production Aurora:**
```hcl
# Terraform configuration
is_aurora = true
engine = "aurora-mysql"
engine_version = "8.0.mysql_aurora.3.03.1"
instance_class = "db.r5.large"
instances = {
  1 = { instance_class = "db.r5.large", promotion_tier = 1 }
  2 = { instance_class = "db.r5.large", promotion_tier = 2 }
}
```

## Environment Configuration

### Environment Variables

**Common Environment Variables:**
```bash
# Database
DB_HOST=database-endpoint
DB_PORT=3306
DB_USER=admin
DB_PASSWORD=secure-password
DB_NAME=sarcng

# Application
PORT=8080
GIN_MODE=release
ENVIRONMENT=production

# Security
JWT_SECRET=your-jwt-secret
```

### Configuration Management

**AWS Secrets Manager:**
```bash
# Store database credentials
aws secretsmanager create-secret \
  --name "sarc-ng/database/credentials" \
  --description "Database credentials for SARC-NG" \
  --secret-string '{"username":"admin","password":"secure-password"}'
```

**AWS Systems Manager Parameter Store:**
```bash
# Store application configuration
aws ssm put-parameter \
  --name "/sarc-ng/config/database-endpoint" \
  --value "your-rds-endpoint.amazonaws.com" \
  --type "String"
```

## Monitoring and Logging

### CloudWatch Integration

**Log Groups:**
- `/ecs/sarc-ng/api` - Application logs
- `/aws/lambda/sarc-ng-api` - Lambda function logs
- `/aws/rds/instance/sarc-ng-db/error` - Database error logs

**Custom Metrics:**
- API request counts and latency
- Database connection pool metrics
- Business logic metrics (reservations, resources)

### Application Observability

**Prometheus Metrics:**
```
# Exposed at /metrics endpoint
http_requests_total
http_request_duration_seconds
database_connections_active
reservation_operations_total
```

**Health Checks:**
```
# Health endpoint
GET /health

# Response
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## SSL/TLS Configuration

### Certificate Management

**AWS Certificate Manager:**
```bash
# Request certificate
aws acm request-certificate \
  --domain-name api.yourdomain.com \
  --validation-method DNS
```

**Load Balancer SSL:**
```hcl
# Terraform ALB listener configuration
port = "443"
protocol = "HTTPS"
ssl_policy = "ELBSecurityPolicy-TLS-1-2-2017-01"
certificate_arn = var.acm_certificate_arn
```

## Domain and DNS

### Route53 Configuration

**DNS Records:**
```hcl
# Terraform Route53 record
resource "aws_route53_record" "api" {
  zone_id = var.hosted_zone_id
  name    = "api.yourdomain.com"
  type    = "A"
  
  alias {
    name                   = aws_lb.main.dns_name
    zone_id                = aws_lb.main.zone_id
    evaluate_target_health = true
  }
}
```

## Security Considerations

### Network Security

**Security Groups:**
- ALB: Allow 80/443 from internet
- ECS Tasks: Allow ALB traffic on app port
- RDS: Allow ECS traffic on 3306

**VPC Configuration:**
- Public subnets for load balancers
- Private subnets for application servers
- Database subnets isolated from internet

### Application Security

**Environment Variables:**
- Use AWS Secrets Manager for sensitive data
- Rotate credentials regularly
- Use IAM roles instead of access keys

**Image Security:**
- Use distroless base images
- Scan images for vulnerabilities
- Use specific image tags, not latest

## Backup and Disaster Recovery

### Database Backups

**Automated Backups:**
- RDS automated backups with point-in-time recovery
- Aurora continuous backups
- Cross-region backup replication for production

**Manual Backups:**
```bash
# Create RDS snapshot
aws rds create-db-snapshot \
  --db-instance-identifier sarc-ng-db \
  --db-snapshot-identifier sarc-ng-backup-$(date +%Y%m%d)
```

### Application Backup

**Container Images:**
- Store in multiple registries
- Tag images with semantic versions
- Keep previous versions for rollback

## Rollback Procedures

### ECS Rollback

**Service Rollback:**
```bash
# Update service to previous task definition
aws ecs update-service \
  --cluster sarc-ng-cluster \
  --service sarc-ng-api \
  --task-definition sarc-ng-api:previous-revision
```

### Lambda Rollback

**Function Rollback:**
```bash
# Update alias to previous version
aws lambda update-alias \
  --function-name sarc-ng-api \
  --name LIVE \
  --function-version previous-version
```

## Troubleshooting

### Common Deployment Issues

**Database Connection Failures:**
```bash
# Check security groups
aws ec2 describe-security-groups --group-ids sg-xxxxx

# Test database connectivity
mysql -h your-endpoint -u admin -p
```

**Container Startup Issues:**
```bash
# Check ECS service events
aws ecs describe-services \
  --cluster sarc-ng-cluster \
  --services sarc-ng-api

# View container logs
aws logs get-log-events \
  --log-group-name /ecs/sarc-ng/api \
  --log-stream-name container/api/task-id
```

**Load Balancer Health Check Failures:**
```bash
# Check target group health
aws elbv2 describe-target-health \
  --target-group-arn your-target-group-arn

# Verify health check endpoint
curl http://your-alb-endpoint/health
```

### Performance Optimization

**ECS Optimization:**
- Right-size CPU and memory allocation
- Use Fargate Spot for cost optimization
- Implement auto-scaling policies

**Database Optimization:**
- Monitor slow query logs
- Optimize connection pool settings
- Use read replicas for read-heavy workloads

**Monitoring Commands:**
```bash
# ECS service metrics
aws cloudwatch get-metric-statistics \
  --namespace AWS/ECS \
  --metric-name CPUUtilization \
  --dimensions Name=ServiceName,Value=sarc-ng-api

# RDS performance insights
aws rds describe-db-instances \
  --db-instance-identifier sarc-ng-db
``` 