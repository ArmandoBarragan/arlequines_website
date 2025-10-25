output "lambda_function_arn" {
  description = "ARN of the Payment Email Lambda Function"
  value       = aws_lambda_function.payment_email.arn
}

output "lambda_function_name" {
  description = "Name of the Payment Email Lambda Function"
  value       = aws_lambda_function.payment_email.function_name
}

output "lambda_deployment_bucket_name" {
  description = "Name of the S3 bucket used for Lambda deployments"
  value       = aws_s3_bucket.lambda_deployments.id
}

output "lambda_role_arn" {
  description = "ARN of the IAM role for the Lambda function"
  value       = aws_iam_role.lambda_role.arn
}

