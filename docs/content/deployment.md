---
sidebar_position: 5
tags:
  - deployment
  - cloud
  - docker
  - aws
---

# Deployment

This guide covers building and deploying SARC-NG in various environments.

## Building for Production

### Local Build

```bash
# Build all binaries (includes Wire code generation)
make build

# Build production release binaries
make release

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

### Docker Build

```bash
# Build Docker image
docker build -t sarc-ng:latest .

# Build with specific tag
docker build -t sarc-ng:v1.0.0 .

# Multi-stage build for smaller production image
docker build --target production -t sarc-ng:prod .
```

### Production Optimizations

```bash
# Build with optimizations
CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -X main.version=${VERSION}" \
  -a -installsuffix cgo \
  -o bin/server cmd/server/main.go

# Strip debugging information
strip bin/server
```

## Local Deployment

### Docker Compose

```bash
# Production-like local deployment
docker compose -f docker-compose.prod.yml up -d

# Scale services
docker compose -f docker-compose.prod.yml up -d --scale sarc-ng=3

# Update specific service
docker compose -f docker-compose.prod.yml up -d --no-deps sarc-ng
```

### Environment Configuration

Create production configuration:

```yaml
# configs/production.yaml
server:
  port: 8080
  timeout: 30s
  graceful_timeout: 10s

database:
  host: ${DB_HOST}
  port: 3306
  name: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  max_connections: 100
  max_idle_connections: 10

auth:
  jwt_secret: ${JWT_SECRET}
  token_ttl: 24h

logging:
  level: info
  format: json
  output: stdout

metrics:
  enabled: true
  port: 9090
  path: /metrics
```

## Cloud Deployment

### AWS Deployment

#### Prerequisites

```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure AWS credentials
aws configure

# Install AWS SAM CLI (for serverless deployment)
pip install aws-sam-cli
# or on macOS
brew install aws-sam-cli

# Install Terraform (for infrastructure as code)
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/
```

#### AWS SAM (Serverless) Deployment

AWS SAM deployment includes Lambda function, RDS MySQL, API Gateway, and VPC configuration.

**Architecture:**

- Lambda Function (Go provided.al2 runtime)
- RDS MySQL 8.0 Database
- API Gateway (HTTP API)
- VPC with public, private, and database subnets
- CloudWatch Logs

**Build Lambda:**

```bash
cd infrastructure/sam
sam build
```

**First-time Deploy (Interactive):**

```bash
cd infrastructure/sam
sam deploy --guided
```

This will prompt for:

- Stack name (e.g., `sarc-ng-prod`)
- AWS Region (e.g., `us-east-1`)
- Environment parameter (dev/staging/prod)
- Database password
- Confirmation for capabilities

**Subsequent Deploys:**

```bash
cd infrastructure/sam
sam deploy
```

**Environment-Specific Deployment:**

```bash
# Development
cd infrastructure/sam && sam deploy --config-env development

# Production
cd infrastructure/sam && sam deploy --config-env production
```

**Custom Parameters:**

```bash
cd infrastructure/sam && sam deploy \
    --parameter-overrides \
    Environment=prod \
    DBPassword=secure-password
```

**Local Testing:**

```bash
# Start database first
make docker-up

# Start SAM local API
cd infrastructure/sam
sam local start-api \
  --port 3001 \
  --docker-network sarc-ng-network \
  --env-vars env.json
```

**Delete Stack:**

```bash
cd infrastructure/sam
sam delete
```

**View Logs:**

```bash
# View Lambda logs
sam logs -n SarcNgFunction --stack-name sarc-ng-prod --tail

# View logs from specific time
sam logs -n SarcNgFunction --stack-name sarc-ng-prod --start-time '10min ago'
```

For detailed SAM documentation, see [`infrastructure/README.md`](../infrastructure/README.md).

#### ECS Deployment

```bash
# Build and push to ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-west-2.amazonaws.com

docker build -t sarc-ng .
docker tag sarc-ng:latest <account-id>.dkr.ecr.us-west-2.amazonaws.com/sarc-ng:latest
docker push <account-id>.dkr.ecr.us-west-2.amazonaws.com/sarc-ng:latest

# Update ECS service
aws ecs update-service --cluster sarc-ng-cluster --service sarc-ng-service --force-new-deployment
```

#### Lambda Deployment

**Recommended:** Use AWS SAM for Lambda deployment (see section above).

**Manual deployment (if not using SAM):**

```bash
# Build Lambda function
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false -o bootstrap cmd/lambda/main.go
zip lambda-deployment.zip bootstrap

# Deploy via AWS CLI
aws lambda update-function-code \
  --function-name sarc-ng-lambda \
  --zip-file fileb://lambda-deployment.zip
```

### Container Platforms

#### Kubernetes

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sarc-ng
  labels:
    app: sarc-ng
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sarc-ng
  template:
    metadata:
      labels:
        app: sarc-ng
    spec:
      containers:
      - name: sarc-ng
        image: sarc-ng:latest
        ports:
        - containerPort: 8080
        env:
        - name: SARC_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: sarc-secrets
              key: db-host
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Scale deployment
kubectl scale deployment sarc-ng --replicas=5

# Rolling update
kubectl set image deployment/sarc-ng sarc-ng=sarc-ng:v2.0.0
```

#### Docker Swarm

```yaml
# docker-stack.yml
version: '3.8'
services:
  sarc-ng:
    image: sarc-ng:latest
    ports:
      - "8080:8080"
    environment:
      - SARC_ENV=production
    networks:
      - sarc-network
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  sarc-network:
    driver: overlay
```

```bash
# Deploy stack
docker stack deploy -c docker-stack.yml sarc-ng

# Update service
docker service update --image sarc-ng:v2.0.0 sarc-ng_sarc-ng
```

## Environment Variables

### Required Environment Variables

```bash
# Database (MySQL)
export DB_HOST=localhost
export DB_PORT=3306
export DB_NAME=sarcng
export DB_USER=root
export DB_PASSWORD=secure_password

# Authentication
export JWT_SECRET=very_secure_jwt_secret

# Server
export SERVER_PORT=8080
export SARC_ENV=production
export GIN_MODE=release
```

### Optional Environment Variables

```bash
# Logging
export LOG_LEVEL=info
export LOG_FORMAT=json

# Metrics
export METRICS_ENABLED=true
export METRICS_PORT=9090

# Features
export FEATURES_RATE_LIMITING=true
export FEATURES_CORS_ENABLED=true

# External Services
export AWS_REGION=us-west-2
export S3_BUCKET=sarc-ng-assets
```

## Database Deployment

### MySQL 8.0 Setup

#### Managed Database (AWS RDS)

**Recommended:** Use AWS SAM deployment which includes RDS MySQL automatically (see SAM section above).

**Manual RDS creation:**

```bash
# Create RDS instance via Terraform
cd infrastructure/terraform
terraform apply -target=aws_db_instance.sarc_ng

# Or via AWS CLI
aws rds create-db-instance \
  --db-instance-identifier sarc-ng-prod \
  --db-instance-class db.t3.micro \
  --engine mysql \
  --engine-version 8.0 \
  --master-username root \
  --master-user-password ${DB_PASSWORD} \
  --db-name sarcng \
  --allocated-storage 20 \
  --storage-type gp2 \
  --vpc-security-group-ids sg-12345678 \
  --db-subnet-group-name sarc-ng-subnet-group
```

#### Self-managed MySQL

```bash
# Run MySQL with Docker
docker run -d \
  --name sarc-mysql \
  -e MYSQL_ROOT_PASSWORD=${DB_PASSWORD} \
  -e MYSQL_DATABASE=sarcng \
  -v mysql-data:/var/lib/mysql \
  -p 3306:3306 \
  mysql:8.0
```

### Database Migrations

**Note:** GORM handles migrations automatically on application startup. For manual control:

```bash
# Run migrations in production
docker run --rm \
  -e DB_HOST=${DB_HOST} \
  -e DB_PORT=3306 \
  -e DB_NAME=sarcng \
  -e DB_USER=root \
  -e DB_PASSWORD=${DB_PASSWORD} \
  sarc-ng:latest \
  /app/cli migrate up

# Or using Kubernetes job
kubectl create job --from=deployment/sarc-ng migrate-job -- /app/cli migrate up
```

## Monitoring Setup

### Prometheus & Grafana

```yaml
# docker-compose.monitoring.yml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
```

### Health Checks

```bash
# Application health endpoint
curl http://your-domain.com/health

# Detailed health check
curl http://your-domain.com/health/detailed

# Metrics endpoint
curl http://your-domain.com/metrics
```

## Load Balancing

### Nginx Configuration

```nginx
# nginx.conf
upstream sarc_ng {
    server sarc-ng-1:8080;
    server sarc-ng-2:8080;
    server sarc-ng-3:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://sarc_ng;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /health {
        access_log off;
        proxy_pass http://sarc_ng/health;
    }
}
```

### AWS Application Load Balancer

```hcl
# Infrastructure code (Terraform)
resource "aws_lb" "sarc_ng" {
  name               = "sarc-ng-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = true
}

resource "aws_lb_target_group" "sarc_ng" {
  name     = "sarc-ng-targets"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = var.vpc_id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/health"
    matcher             = "200"
  }
}
```

## SSL/TLS Configuration

### Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### AWS Certificate Manager

```hcl
resource "aws_acm_certificate" "sarc_ng" {
  domain_name       = "your-domain.com"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}
```

## Deployment Strategies

### Blue-Green Deployment

```bash
# Deploy to green environment
docker service create --name sarc-ng-green sarc-ng:v2.0.0

# Test green environment
curl http://green.your-domain.com/health

# Switch traffic (update load balancer)
# Remove blue environment after validation
docker service rm sarc-ng-blue
```

### Rolling Updates

```bash
# Kubernetes rolling update
kubectl set image deployment/sarc-ng sarc-ng=sarc-ng:v2.0.0

# Docker Swarm rolling update
docker service update --image sarc-ng:v2.0.0 sarc-ng_sarc-ng
```

### Canary Deployment

```bash
# Deploy canary version (10% traffic)
kubectl patch deployment sarc-ng-canary -p '{"spec":{"replicas":1}}'
kubectl patch deployment sarc-ng-stable -p '{"spec":{"replicas":9}}'

# Monitor metrics, gradually increase canary traffic
# Rollback if issues detected
kubectl patch deployment sarc-ng-canary -p '{"spec":{"replicas":0}}'
```

## Troubleshooting Deployment

### Common Issues

**Container Won't Start**

```bash
# Check container logs
docker logs sarc-ng-container

# Check resource constraints
docker stats sarc-ng-container

# Exec into container
docker exec -it sarc-ng-container /bin/sh
```

**Database Connection Issues**

```bash
# Test database connectivity
docker run --rm -it postgres:15 psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Check network connectivity
nc -zv $DB_HOST $DB_PORT
```

**High Memory Usage**

```bash
# Profile memory usage
curl http://your-domain.com/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Adjust container limits
docker update --memory=1g sarc-ng-container
```

**SSL Certificate Issues**

```bash
# Check certificate validity
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# Renew Let's Encrypt certificate
sudo certbot renew
```

This deployment guide provides comprehensive coverage for getting SARC-NG running in production environments with proper security, monitoring, and scalability considerations.
