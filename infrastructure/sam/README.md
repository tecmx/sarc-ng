# SARC-NG AWS SAM Deployment

This directory contains AWS SAM (Serverless Application Model) templates for deploying SARC-NG as a serverless application using AWS Lambda and RDS.

## Prerequisites

1. AWS CLI configured with appropriate credentials
2. AWS SAM CLI installed
3. Go 1.x installed
4. Access to AWS Academy LabRole

## Quick Start

1. **Build the Lambda function:**
   ```bash
   ./build.sh
   ```

2. **Build the SAM application:**
   ```bash
   sam build
   ```

3. **Deploy with guided setup (first time):**
   ```bash
   sam deploy --guided
   ```

4. **Deploy with existing configuration:**
   ```bash
   sam deploy
   ```

## Architecture

The SAM template deploys:

- **AWS Lambda Function**: Runs the SARC-NG Go application
- **Amazon RDS MySQL**: Database for the application
- **API Gateway**: HTTP API endpoint for the Lambda function
- **CloudWatch Logs**: Log groups for monitoring

## Configuration

### Parameters

- `Environment`: Deployment environment (dev, staging, prod)
- `DBPassword`: Password for the RDS MySQL instance
- `DBName`: Name of the database to create

### Environment Variables

The Lambda function receives these environment variables:
- `DB_HOST`: RDS endpoint
- `DB_PORT`: Database port (3306)
- `DB_NAME`: Database name
- `DB_USER`: Database username (root)
- `DB_PASSWORD`: Database password
- `GIN_MODE`: Gin framework mode (release)
- `LOG_LEVEL`: Application log level (info)
- `ENVIRONMENT`: Deployment environment

## Local Development

1. **Start local API:**
   ```bash
   sam local start-api
   ```

2. **Invoke function locally:**
   ```bash
   sam local invoke SarcNgFunction
   ```

3. **Sync changes (watch mode):**
   ```bash
   sam sync --watch
   ```

## Deployment Environments

### Development
```bash
sam deploy --config-env development
```

### Production
```bash
sam deploy --config-env production
```

## Commands Reference

### Build
```bash
# Build SAM application
sam build

# Build with cache
sam build --cached

# Build in parallel
sam build --parallel
```

### Deploy
```bash
# Guided deployment (first time)
sam deploy --guided

# Deploy with specific parameters
sam deploy --parameter-overrides Environment=dev DBPassword=mypassword

# Deploy to specific environment
sam deploy --config-env production
```

### Local Testing
```bash
# Start local API Gateway
sam local start-api --port 3000

# Invoke specific function
sam local invoke SarcNgFunction --event events/api-gateway.json

# Generate sample events
sam local generate-event apigateway aws-proxy > events/api-gateway.json
```

### Cleanup
```bash
# Delete the stack
sam delete

# Delete with specific stack name
sam delete --stack-name sarc-ng-dev
```

## Monitoring

After deployment, monitor your application using:

1. **CloudWatch Logs**: Check Lambda function logs
2. **CloudWatch Metrics**: Monitor Lambda performance
3. **API Gateway Logs**: Check API access logs
4. **RDS Monitoring**: Database performance metrics

## Troubleshooting

### Common Issues

1. **Lambda timeout**: Increase the timeout in the Globals section
2. **Database connection**: Check security groups and network configuration
3. **Build failures**: Ensure Go binary is built for Linux AMD64

### Useful Commands

```bash
# Check SAM version
sam --version

# Validate template
sam validate

# Get stack outputs
aws cloudformation describe-stacks --stack-name sarc-ng-dev --query 'Stacks[0].Outputs'

# Check function logs
sam logs -n SarcNgFunction --stack-name sarc-ng-dev --tail
```

## Security Notes

- The template uses AWS Academy LabRole for simplicity
- RDS instance is publicly accessible for development
- In production, consider using VPC with private subnets
- Store sensitive parameters in AWS Parameter Store or Secrets Manager

## Cost Optimization

- Lambda: Pay per request and execution time
- RDS: db.t3.micro for development (eligible for free tier)
- API Gateway: Pay per API call
- CloudWatch: Basic monitoring included

## Further Reading

- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [AWS Lambda Go Documentation](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)
- [API Gateway Documentation](https://docs.aws.amazon.com/apigateway/)
