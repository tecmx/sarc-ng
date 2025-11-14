# Lambda SQS Consumer Module Example

Creates a Lambda function that processes messages from an SQS queue.

## Prerequisites

- Existing SQS queue
- Lambda deployment package

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Lambda Function
- SQS Event Source Mapping
- Dead Letter Queue (optional)
- CloudWatch Log Group
- IAM Role and Policies

## Cleanup

```bash
terraform destroy
```

