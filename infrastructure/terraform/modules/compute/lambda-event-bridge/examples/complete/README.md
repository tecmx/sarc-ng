# Lambda EventBridge Module Example

Creates a Lambda function triggered by EventBridge (scheduled or event-pattern based).

## Prerequisites

- Lambda deployment package
- AWS account with permissions

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Resources Created

- Lambda Function
- EventBridge Rule
- EventBridge Target
- CloudWatch Log Group
- IAM Role and Policies
- Lambda Permission for EventBridge

## Cleanup

```bash
terraform destroy
```

