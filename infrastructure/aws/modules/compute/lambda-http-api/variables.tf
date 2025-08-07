variable "name" {
  description = "Name of the Lambda function"
  type        = string
}

variable "description" {
  description = "Description of the Lambda function"
  type        = string
  default     = ""
}

variable "handler" {
  description = "Lambda function handler"
  type        = string
}

variable "runtime" {
  description = "Lambda function runtime"
  type        = string
}

variable "memory_size" {
  description = "Amount of memory in MB the Lambda Function can use at runtime"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Amount of time the Lambda Function has to run in seconds"
  type        = number
  default     = 3
}

variable "source_path" {
  description = "Path to the Lambda function source code"
  type        = string
  default     = null
}

variable "source_bucket" {
  description = "S3 bucket containing the Lambda deployment package"
  type        = string
  default     = null
}

variable "source_key" {
  description = "S3 key of the Lambda deployment package"
  type        = string
  default     = null
}

variable "environment_variables" {
  description = "Environment variables for the Lambda function"
  type        = map(string)
  default     = {}
}

variable "vpc_config" {
  description = "VPC configuration for the Lambda function"
  type = object({
    subnet_ids         = list(string)
    security_group_ids = list(string)
  })
  default = null
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}

variable "api_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = null
}

variable "api_description" {
  description = "Description of the API Gateway"
  type        = string
  default     = null
}

variable "api_path" {
  description = "Path part for the API Gateway resource"
  type        = string
  default     = "{proxy+}"
}

variable "api_stage_name" {
  description = "Name of the API Gateway stage"
  type        = string
  default     = "api"
}

variable "api_domain_name" {
  description = "Custom domain name for the API Gateway"
  type        = string
  default     = null
}

variable "api_certificate_arn" {
  description = "ARN of the ACM certificate for the custom domain"
  type        = string
  default     = null
}

variable "cors_configuration" {
  description = "CORS configuration for the API Gateway"
  type = object({
    allow_origins     = list(string)
    allow_methods     = list(string)
    allow_headers     = list(string)
    expose_headers    = list(string)
    allow_credentials = bool
    max_age           = number
  })
  default = null
}

variable "log_retention_in_days" {
  description = "Number of days to retain CloudWatch logs"
  type        = number
  default     = 14
}

variable "iam_role_arn" {
  description = "ARN of an existing IAM role for the Lambda function"
  type        = string
  default     = null
}

variable "iam_policy_documents" {
  description = "List of IAM policy documents to attach to the Lambda role"
  type        = list(string)
  default     = []
}

variable "create_api_gateway" {
  description = "Whether to create an API Gateway for the Lambda function"
  type        = bool
  default     = true
}

variable "create_custom_domain" {
  description = "Whether to create a custom domain for the API Gateway"
  type        = bool
  default     = false
}

variable "enable_xray_tracing" {
  description = "Whether to enable X-Ray tracing for the Lambda function"
  type        = bool
  default     = false
}

variable "reserved_concurrent_executions" {
  description = "Amount of reserved concurrent executions for the Lambda function"
  type        = number
  default     = -1
}

variable "publish" {
  description = "Whether to publish creation/change as a new Lambda function version"
  type        = bool
  default     = false
}

variable "create_async_event_config" {
  description = "Whether to create async event configuration for the Lambda function"
  type        = bool
  default     = false
}

variable "maximum_retry_attempts" {
  description = "Maximum number of retry attempts for async invocations"
  type        = number
  default     = 2
}

variable "maximum_event_age_in_seconds" {
  description = "Maximum age of a request that Lambda sends to a function for processing"
  type        = number
  default     = 60
}

variable "destination_on_failure" {
  description = "ARN of the destination resource for failed async invocations"
  type        = string
  default     = null
}

variable "destination_on_success" {
  description = "ARN of the destination resource for successful async invocations"
  type        = string
  default     = null
} 
