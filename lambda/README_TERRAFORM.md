# Terraform Deployment for Lambda Function

This directory contains Terraform configuration to deploy the Payment Email Lambda function using S3 for deployment packages.

## Prerequisites

1. **Terraform installed** (>= 1.0)
2. **AWS CLI configured** with appropriate credentials
3. **Go installed** (for building the Lambda function)

## Setup

1. **Copy the example variables file:**
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. **Edit `terraform.tfvars`** with your values:
   - Set `smtp_pass` to your SMTP password
   - Adjust other variables as needed

## Deployment

### Using Make (Recommended)

The Makefile handles building and deploying:

```bash
export SMTP_PASSWORD=your_password
export AWS_REGION=us-east-1
make deploy-terraform
```

### Manual Deployment

1. **Build the Lambda function:**
   ```bash
   make build-terraform
   ```

2. **Initialize Terraform:**
   ```bash
   terraform init
   ```

3. **Plan the deployment:**
   ```bash
   terraform plan -var="smtp_pass=your_password"
   ```

4. **Apply the configuration:**
   ```bash
   terraform apply -var="smtp_pass=your_password"
   ```

## How It Works

1. **Build**: The Go binary is compiled for Linux (`GOOS=linux GOARCH=amd64`) and named `bootstrap` (required for `provided.al2023` runtime)

2. **Zip**: Terraform's `archive_file` data source creates a zip file containing the `bootstrap` binary

3. **Upload**: The zip file is uploaded to an S3 bucket with versioning enabled

4. **Deploy**: The Lambda function is created/updated using the S3 object as the deployment package

5. **Event Source**: An SQS event source mapping connects the Lambda to your SQS queue

## Resources Created

- **S3 Bucket**: Stores Lambda deployment packages (with versioning)
- **IAM Role**: Grants Lambda permissions for CloudWatch Logs and SQS
- **Lambda Function**: The payment email sender function
- **Event Source Mapping**: Connects SQS queue to Lambda
- **CloudWatch Log Group**: For Lambda function logs

## Variables

See `variables.tf` for all available variables. Key variables:

- `smtp_pass`: SMTP password (sensitive, required)
- `aws_region`: AWS region (default: us-east-1)
- `lambda_deployment_bucket_name`: S3 bucket name for deployments
- `sqs_queue_arn`: ARN of the SQS queue to trigger Lambda

## Outputs

After deployment, you can view outputs:

```bash
terraform output
```

Available outputs:
- `lambda_function_arn`: ARN of the Lambda function
- `lambda_function_name`: Name of the Lambda function
- `lambda_deployment_bucket_name`: S3 bucket used for deployments
- `lambda_role_arn`: IAM role ARN

## Cleanup

To destroy all resources:

```bash
terraform destroy -var="smtp_pass=your_password"
```

## Notes

- The `bootstrap` binary must exist before running `terraform apply`. Use `make build-terraform` or `make deploy-terraform` to ensure it's built.
- The S3 bucket name must be globally unique. Change `lambda_deployment_bucket_name` if needed.
- The zip file is automatically created by Terraform and doesn't need to be committed to git.

