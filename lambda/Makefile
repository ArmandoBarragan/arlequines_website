# Build the Lambda binary for deployment
build:
	@echo "Building Lambda function binary..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	@echo "Converting binary to zip file..."
	@zip -j lambda-function.zip bootstrap
	@echo "Zip file created successfully. Build complete."

# Check required environment variables
check-env:
	@if [ -z "$$SMTP_PASSWORD" ] || [ -z "$$AWS_REGION" ]; then \
		echo "Error: SMTP_PASSWORD and AWS_REGION environment variables must be set"; \
		echo "To fix this, run: export SMTP_PASSWORD=your_password && export AWS_REGION=your_region"; \
		echo "Or pass it inline: SMTP_PASSWORD=your_password AWS_REGION=your_region make deploy"; \
		exit 1; \
	fi

# Deploy Lambda function using Terraform (S3-based)
deploy: build check-env
	@echo "Verifying bootstrap file is not empty..."
	@if [ ! -s lambda-function.zip ]; then \
		echo "Error: zip file is empty."; \
		exit 1; \
	fi
	@echo "Zip file verified (size: $$(ls -lh lambda-function.zip | awk '{print $$5}'))"
	@echo "Deploying with Terraform..."
	@terraform init
	@terraform plan -var="smtp_pass=$$SMTP_PASSWORD" -var="aws_region=$$AWS_REGION"
	@terraform apply -var="smtp_pass=$$SMTP_PASSWORD" -var="aws_region=$$AWS_REGION" -auto-approve
	@echo "Lambda function deployed successfully"

# Check SMTP_PASSWORD only (for local testing)
check-smtp:
	@if [ -z "$$SMTP_PASSWORD" ]; then \
		echo "Error: SMTP_PASSWORD environment variable is not set or not exported"; \
		echo "To fix this, run: export SMTP_PASSWORD=your_password"; \
		echo "Or pass it inline: SMTP_PASSWORD=your_password make test-local"; \
		exit 1; \
	fi

# Test Lambda function locally (uses SAM CLI for local testing)
# Note: SAM CLI is used only for local testing, deployment uses Terraform
test-local: build check-smtp
	@echo "Testing Lambda function locally"
	@sam local invoke --event test-event.json --parameter-overrides "SMTPPass=$$SMTP_PASSWORD"
	@echo "Lambda function tested successfully"
