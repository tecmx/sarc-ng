# SARC-NG SAM Guide

> Complete guide for local development and AWS deployment with SAM

## Quick Start

### Local Development (No AWS Credentials Needed)

```bash
cd infrastructure/sam

# Start everything (MySQL + API Gateway)
make local-dev

# Stop everything
make local-stop
```

**Configuration:** Edit `local.env` for environment variables (simple key=value format).
`env.json` is auto-generated from `local.env` automatically.

### Deploy to AWS Production

```bash
# Install SAM CLI (if needed)
wget https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip
unzip aws-sam-cli-linux-x86_64.zip -d sam-installation
sudo ./sam-installation/install

# Deploy everything
make deploy-all ENV=prod

# See endpoints
make urls ENV=prod
```

---

## Stack Architecture

```
📦 sarc-ng-cognito-{env}          📦 sarc-ng-{env}
├── Cognito User Pool              ├── VPC & Networking
├── User Pool Client               ├── RDS Database
├── User Groups                    ├── Lambda Function
├── SSM Parameters                 ├── API Gateway
└── CloudFormation Exports     ◄───┤ Custom Resource
                                    └── (Auto-adds Swagger callback)
```

## Files

- **`cognito-stack.yaml`** - Authentication infrastructure (Cognito)
- **`api-stack.yaml`** - API infrastructure (Lambda, RDS, VPC)
- **`Makefile`** - Deployment automation
- **`local.env`** - Local development configuration

---

## Local Development

### Prerequisites

- Docker and Docker Compose
- AWS SAM CLI (`pip install aws-sam-cli`)
- **No AWS credentials required** for local development

### Commands

| What | Command |
|------|---------|
| **Full local environment** | `make local-dev` |
| **Stop local SAM** | `make local-stop` |
| **Just Docker services** | `cd ../docker && docker compose up -d` |
| **Stop Docker services** | `cd ../docker && docker compose down` |

### Local Endpoints

- API: <http://localhost:3000/api/v1>
- Swagger: <http://localhost:3000/swagger/index.html>
- Database: `localhost:3306` (user: `root`, password: from `local.env`)

### Troubleshooting Local

**Validate Template**

```bash
sam validate
```

**Debug Invoke**

```bash
sam local invoke SarcNgFunction --event events/test.json --debug
```

**Check Build**

```bash
ls -la .aws-sam/build/SarcNgFunction/
```

---

## AWS Deployment

### Prerequisites

- AWS CLI with credentials configured
- AWS SAM CLI installed

### Deploy Commands

```bash
# Deploy to specific environment
make deploy-all ENV=dev
make deploy-all ENV=staging
make deploy-all ENV=prod

# Deploy only Cognito
make deploy-cognito ENV=dev

# Deploy only API
make deploy-api ENV=dev
```

### Manage

```bash
# Check status
make status ENV=prod

# Get URLs
make urls ENV=prod

# Validate templates
make validate

# Clean build artifacts
make clean
```

### Delete

```bash
# Delete API stack
make delete-api ENV=dev

# Delete Cognito (⚠️ deletes users!)
make delete-cognito ENV=dev

# Delete everything
make delete-all ENV=dev
```

### Environment Variables

Required in `../../aws.env`:

```bash
# AWS Credentials
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
AWS_DEFAULT_REGION=us-east-1

# Database (optional - can be set during deploy)
DB_PASSWORD=YourSecurePassword
```

Load before deployment:

```bash
source ../../aws.env
```

---

## First AWS Deployment

```bash
# 1. Load AWS credentials
cd infrastructure/sam
source ../../aws.env

# 2. Validate templates
make validate

# 3. Deploy both stacks
make deploy-all ENV=prod

# 4. Get URLs
make urls ENV=prod

# 5. Create users
USER_POOL_ID=$(aws cloudformation describe-stacks \
  --stack-name sarc-ng-cognito-prod \
  --query 'Stacks[0].Outputs[?OutputKey==`CognitoUserPoolId`].OutputValue' \
  --output text)

aws cognito-idp admin-create-user \
  --user-pool-id $USER_POOL_ID \
  --username admin \
  --user-attributes Name=email,Value=admin@example.com \
  --temporary-password "Admin123!"

aws cognito-idp admin-add-user-to-group \
  --user-pool-id $USER_POOL_ID \
  --username admin \
  --group-name admin
```

---

## Testing with Swagger

1. **Get Swagger URL:**

   ```bash
   make urls ENV=prod
   ```

2. **Open Swagger UI in browser**

3. **Click "Authorize"** - Swagger will redirect to Cognito Hosted UI

4. **Login with Cognito credentials**

5. **Make authenticated API calls**

---

## Deployment Order

**Critical:** Always deploy in this order:

1. ✅ Deploy **Cognito stack** first
2. ✅ Deploy **API stack** second (imports Cognito)

The API stack depends on Cognito exports. Deploying in wrong order will fail.

---

## What Gets Created

### Cognito Stack

- **User Pool** - Manages users
- **User Pool Client** - OAuth configuration
- **User Pool Domain** - Hosted UI
- **User Groups** - admin, manager, teacher, student
- **SSM Parameters** - Configuration storage
- **CloudFormation Exports** - For API stack import

### API Stack

- **VPC** - Network isolation (3 AZs)
- **RDS MySQL** - Database
- **Secrets Manager** - DB credentials
- **Lambda Function** - API logic
- **API Gateway** - HTTP endpoint
- **Custom Resource** - Adds Swagger callback to Cognito

---

## Troubleshooting

### "No export named ... found"

Deploy Cognito first:

```bash
make deploy-cognito ENV=prod
```

### Stack update failed

Check events:

```bash
aws cloudformation describe-stack-events \
  --stack-name sarc-ng-prod \
  --max-items 20
```

### Custom Resource failed

Check logs:

```bash
aws logs tail /aws/lambda/prod-sarc-ng-update-cognito-callbacks --follow
```

---

## Architecture Benefits

✅ **No circular dependencies**
✅ **Independent stack updates**
✅ **Cleaner infrastructure**
✅ **Swagger UI auto-configured**
✅ **Environment isolation**
✅ **Easy rollbacks**
✅ **Local development without AWS**

---

## Production Checklist

Before deploying to production:

- [ ] Review and update callback URLs
- [ ] Enable MFA (`MFAConfiguration=ON`)
- [ ] Enable Advanced Security (`AdvancedSecurityMode=ENFORCED`)
- [ ] Use strong DB password
- [ ] Review security groups
- [ ] Set up CloudWatch alarms
- [ ] Test in dev/staging first
- [ ] Document user creation process
- [ ] Set up backup strategy for Cognito users
- [ ] Configure custom domain (optional)

---

## Support

For issues, check:

1. CloudFormation Events (see troubleshooting)
2. CloudWatch Logs
3. Stack Outputs (`make urls`)
4. Template validation (`make validate`)
