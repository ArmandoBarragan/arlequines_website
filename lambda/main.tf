terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# S3 bucket for Lambda deployment package
resource "aws_s3_bucket" "lambda_deployments" {
  bucket = var.lambda_deployment_bucket_name

  tags = {
    Name        = "Lambda Deployment Bucket"
    Environment = var.environment
  }
}

resource "aws_s3_bucket_versioning" "lambda_deployments" {
  bucket = aws_s3_bucket.lambda_deployments.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "lambda_deployments" {
  bucket = aws_s3_bucket.lambda_deployments.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# Archive the Lambda function binary
# Note: bootstrap file must exist before running terraform apply
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "${path.module}/lambda.zip"
}

# Upload Lambda deployment package to S3
resource "aws_s3_object" "lambda_zip" {
  bucket = aws_s3_bucket.lambda_deployments.id
  key    = "payment-email-sender/${data.archive_file.lambda_zip.output_md5}.zip"
  source = data.archive_file.lambda_zip.output_path
  etag   = data.archive_file.lambda_zip.output_md5

  tags = {
    Name        = "Payment Email Lambda"
    Environment = var.environment
  }
}

# IAM role for Lambda function
resource "aws_iam_role" "lambda_role" {
  name = "payment-email-sender-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Name        = "Payment Email Lambda Role"
    Environment = var.environment
  }
}

# IAM policy for Lambda to access SQS and CloudWatch
resource "aws_iam_role_policy" "lambda_policy" {
  name = "payment-email-sender-lambda-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ]
        Resource = var.sqs_queue_arn
      }
    ]
  })
}

# Lambda function
resource "aws_lambda_function" "payment_email" {
  function_name = "payment-email-sender"
  role          = aws_iam_role.lambda_role.arn
  handler       = "main"
  runtime       = "provided.al2023"
  timeout       = 30
  memory_size   = 128

  s3_bucket = aws_s3_bucket.lambda_deployments.id
  s3_key    = aws_s3_object.lambda_zip.key

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  environment {
    variables = {
      SMTPHost  = var.smtp_host
      SMTPPort  = var.smtp_port
      SMTPUser  = var.smtp_user
      SMTPPass  = var.smtp_pass
      FromEmail = var.from_email
    }
  }

  tags = {
    Name        = "Payment Email Sender"
    Environment = var.environment
  }
}

# SQS event source mapping
resource "aws_lambda_event_source_mapping" "sqs_trigger" {
  event_source_arn = var.sqs_queue_arn
  function_name    = aws_lambda_function.payment_email.arn
  batch_size       = 10
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.payment_email.function_name}"
  retention_in_days = 14

  tags = {
    Name        = "Payment Email Lambda Logs"
    Environment = var.environment
  }
}

