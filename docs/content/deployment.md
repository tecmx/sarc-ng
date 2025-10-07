---
sidebar_position: 5
tags:
  - deployment
  - ci-cd
  - github-actions
---

# Deployment

## CI/CD Pipeline

SARC-NG uses GitHub Actions for automated building, testing, and deployment.

### Workflows

#### 🔄 CI (Continuous Integration)

**Trigger:** Push or PR to `main` or `develop`

**Jobs:**
1. **Lint** - Code quality checks with golangci-lint
2. **Unit Tests** - Run all unit tests with coverage
3. **Integration Tests** - Test with MySQL database
4. **Build** - Build all binaries
5. **Security** - Scan for vulnerabilities (Gosec, Trivy)

**Status:** [![CI](https://github.com/tecmx/sarc-ng/workflows/CI/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/ci.yml)

#### 🚀 Deploy (Continuous Deployment)

**Trigger:**
- Push to `main` → Deploy to **production**
- Push to `develop` → Deploy to **dev**
- Manual trigger → Choose environment

**Steps:**
1. Build SAM application
2. Deploy to AWS Lambda + RDS
3. Run health checks
4. Post deployment summary

**Status:** [![Deploy](https://github.com/tecmx/sarc-ng/workflows/Deploy/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/deploy.yml)

#### 📦 Release

**Trigger:** Push tag `v*.*.*` (e.g., `v1.0.0`)

**Artifacts:**
- Binaries for Linux (amd64, arm64)
- Binaries for macOS (amd64, arm64)
- Lambda bootstrap binary
- Docker image pushed to GHCR
- Checksums and changelog

**Status:** [![Release](https://github.com/tecmx/sarc-ng/workflows/Release/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/release.yml)

### Setup Requirements

#### GitHub Secrets

Configure these secrets in your repository settings:

```
AWS_ACCESS_KEY_ID       - AWS access key for deployment
AWS_SECRET_ACCESS_KEY   - AWS secret key for deployment
DB_PASSWORD             - Database password for SAM deployment
```

#### Environment Protection (Optional)

Configure branch protection and environment rules:

1. **Settings** → **Environments**
2. Create: `dev`, `staging`, `production`
3. Add reviewers for production
4. Configure branch restrictions

### Usage

#### Manual Deployment

Deploy to specific environment via GitHub Actions:

```yaml
# Go to Actions → Deploy → Run workflow
# Select environment: dev/staging/production
```

#### Creating a Release

```bash
# Tag a new version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Workflow automatically:
# - Builds all binaries
# - Creates GitHub release
# - Builds and pushes Docker image
```

#### Local Testing

Test workflows locally with [act](https://github.com/nektos/act):

```bash
# Install act
brew install act  # macOS

# Run CI workflow
act push -W .github/workflows/ci.yml

# Run specific job
act -j test
```

### Deployment Environments

**Development (`develop` branch)**
- Stack: `sarc-ng-dev`
- Auto-deploy on push
- Lower resources
- Faster iterations

**Staging (manual trigger)**
- Stack: `sarc-ng-staging`
- Manual deployment
- Production-like config
- Pre-release testing

**Production (`main` branch)**
- Stack: `sarc-ng-prod`
- Auto-deploy on push to main
- Full resources
- Protected environment

### Monitoring CI/CD

**Build Status**
- Actions tab in GitHub
- Status badges in README
- Email notifications (Settings)

**Deployment URLs**

After deployment, check job summary for:
- API endpoint URL
- Swagger documentation URL
- Health check status

**Logs**

```bash
# CloudWatch (AWS)
aws logs tail /aws/lambda/sarc-ng-prod --follow

# GitHub Actions
# Actions → Select workflow → View logs
```

### Rollback

**Revert Deployment:**

```bash
# Rollback to previous version
git revert <commit-hash>
git push origin main

# Or redeploy specific tag
git checkout v1.0.0
# Manually trigger Deploy workflow
```

**Manual SAM Rollback:**

```bash
cd infrastructure/sam
sam deploy --config-env production \
  --parameter-overrides Version=v1.0.0
```

## Build Commands

```bash
make build          # Build all binaries
make release        # Production build with optimizations
```

## Docker Deployment

```bash
# Build image
docker build -t sarc-ng:latest .

# Run production stack
docker compose -f docker-compose.prod.yml up -d
```

## AWS SAM (Serverless)

**Prerequisites:** AWS CLI and SAM CLI installed, credentials configured.

### Deploy

```bash
cd infrastructure/sam

# First time (interactive)
sam build && sam deploy --guided

# Subsequent deploys
sam build && sam deploy
```

### Local Testing

```bash
# Start database
make docker-up

# Start local API
cd infrastructure/sam
sam local start-api --port 3001 \
    --docker-network sarc-ng-network \
  --env-vars env.json
```

## AWS Terraform

```bash
cd infrastructure/terraform/live/dev
terraform init
terraform plan
terraform apply
```

## Configuration

Key environment variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_NAME=sarcng
DB_USER=root
DB_PASSWORD=<password>

# Auth
JWT_SECRET=<secret>

# Server
SERVER_PORT=8080
SARC_ENV=production
GIN_MODE=release
```

See `configs/production.yaml` for full configuration options.

## Monitoring

```bash
# Health check
curl http://your-domain.com/health

# Metrics
curl http://your-domain.com/metrics
```

## Troubleshooting

```bash
# Check logs
docker logs <container-name>

# Test database connection
nc -zv $DB_HOST $DB_PORT

# Check container resources
docker stats <container-name>
```
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

## Best Practices

### CI/CD Best Practices

1. **Always create PRs** - Never push directly to main
2. **Wait for CI** - Ensure all checks pass before merging
3. **Tag releases** - Use semantic versioning (v1.0.0)
4. **Test deployments** - Use staging before production
5. **Monitor logs** - Check CloudWatch after deployment
6. **Use secrets** - Never hardcode credentials
7. **Review changes** - Use SAM changeset review before deploy

### Deployment Best Practices

1. **Backup before deployment** - Always backup database
2. **Use blue-green deployments** - Minimize downtime
3. **Monitor after deployment** - Watch logs and metrics
4. **Have rollback plan** - Document rollback procedures
5. **Test in staging first** - Never test in production
6. **Use configuration management** - Keep configs in version control
7. **Document incidents** - Learn from deployment issues

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [Semantic Versioning](https://semver.org/)
- [Twelve-Factor App](https://12factor.net/)

---

This deployment guide provides comprehensive coverage for getting SARC-NG running in production environments with proper CI/CD automation, security, monitoring, and scalability considerations.
