# Lambda HTTP API Module Example

Creates a Lambda function with HTTP API Gateway integration.

## Prerequisites

- Lambda deployment package (function.zip)
- AWS account with Lambda permissions

## Usage

Create a simple Lambda function:

```python
# handler.py
def handler(event, context):
    return {
        'statusCode': 200,
        'body': json.dumps({'message': 'Hello from Lambda!'})
    }
```

Package and deploy:

```bash
# Create deployment package
zip function.zip handler.py

# Deploy
terraform init
terraform plan
terraform apply
```

## Resources Created

- Lambda Function
- HTTP API Gateway v2
- API Gateway Stage
- CloudWatch Log Groups
- IAM Role and Policies
- Lambda Permission for API Gateway

## Cleanup

```bash
terraform destroy
```

