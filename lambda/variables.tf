variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "lambda_deployment_bucket_name" {
  description = "Name of the S3 bucket for Lambda deployments"
  type        = string
  default     = "arlequines-lambda-deployments"
}

variable "smtp_host" {
  description = "SMTP server host"
  type        = string
  default     = "smtp.gmail.com"
}

variable "smtp_port" {
  description = "SMTP server port"
  type        = string
  default     = "587"
}

variable "smtp_user" {
  description = "SMTP username"
  type        = string
  default     = "armandobp765@gmail.com"
}

variable "smtp_pass" {
  description = "SMTP password"
  type        = string
  sensitive   = true
}

variable "from_email" {
  description = "From email address"
  type        = string
  default     = "armandobp765@gmail.com"
}

variable "sqs_queue_arn" {
  description = "ARN of the SQS queue to trigger the Lambda"
  type        = string
  default     = "arn:aws:sqs:us-east-1:962469996968:Arlequines.fifo"
}

