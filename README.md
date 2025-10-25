# Arlequines Website

A theater booking platform built with Go, featuring play and presentation management, payment processing with Stripe, and automated email confirmations via AWS Lambda.

## Features

- **Play Management**: Browse and view available plays with author information
- **Presentation Management**: View upcoming presentations sorted by date
- **Payment Processing**: Integrated Stripe payment processing for ticket purchases
- **User Authentication**: JWT-based authentication system
- **Admin Panel**: Administrative routes for managing content
- **Email Notifications**: Automated payment confirmation emails via AWS Lambda and SQS
- **Docker Support**: Easy local development with Docker Compose

## Tech Stack

- **Backend**: Go 1.23 with [Fiber](https://gofiber.io/) web framework
- **Database**: PostgreSQL 15 with [GORM](https://gorm.io/)
- **Authentication**: JWT tokens (golang-jwt/jwt/v5)
- **Payment**: [Stripe](https://stripe.com/) integration
- **Caching**: Redis (optional, currently commented out)
- **Cloud**: AWS Lambda for email processing, SQS for message queuing
- **Infrastructure**: Terraform for Lambda deployment
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL 15 (or use Docker Compose)
- AWS CLI (for Lambda deployment)
- Terraform (for Lambda deployment)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd arlequines_website
```

### 2. Environment Configuration

Copy the example environment file and configure it:

```bash
cp example.env .env
```

Edit `.env` with your configuration:

```env
# Database
POSTGRES_USER=arlequines
POSTGRES_PASSWORD=your_password
POSTGRES_DB=arlequines_db
DB_HOST=db
DB_PORT=5432
DB_USER=arlequines
DB_PASSWORD=your_password
DB_NAME=arlequines_db
DB_SSLMODE=disable

# Application
HOST_URL=http://localhost:8000
SECRET_KEY=your_secret_key_for_jwt

# Stripe
STRIPE_PUBLIC_KEY=your_stripe_public_key
STRIPE_PRIVATE_KEY=your_stripe_private_key

# AWS (for Lambda email service)
AWS_REGION=us-east-1
SQS_QUEUE_URL=your_sqs_queue_url
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key

# Redis (optional)
REDIS_PASSWORD=your_redis_password
```

### 3. Run with Docker Compose

Start the application and database:

```bash
docker-compose up -d
```

The application will be available at `http://localhost:8000`.

### 4. Run Locally (without Docker)

If you prefer to run the application directly:

```bash
# Start PostgreSQL (or use Docker Compose for just the database)
docker-compose up -d db

# Run the application
cd app
go run main.go
```

## Project Structure

```
arlequines_website/
├── app/                    # Main application
│   ├── main.go            # Application entry point
│   ├── Dockerfile         # Docker configuration
│   ├── go.mod             # Go dependencies
│   ├── settings/          # Configuration and database setup
│   └── src/
│       ├── handlers/      # HTTP request handlers
│       ├── models/        # Database models (Play, Presentation, User)
│       ├── repositories/  # Data access layer
│       ├── routers/       # Route definitions
│       └── services/      # Business logic
├── lambda/                # AWS Lambda function for email processing
│   ├── main.go           # Lambda handler
│   ├── main.tf           # Terraform configuration
│   ├── Makefile          # Build and deployment scripts
│   └── README_TERRAFORM.md # Lambda deployment guide
├── docker-compose.yml     # Docker Compose configuration
├── example.env            # Example environment variables
└── README.md             # This file
```

## API Endpoints

### Public Routes

- `GET /plays` - List all plays (sorted alphabetically)
- `GET /plays/:id` - Get play details
- `GET /presentations` - List all presentations (sorted by date)
- `GET /presentations/:id` - Get presentation details

### Authentication Routes

- Authentication endpoints (see `app/src/routers/auth.go`)

### Stripe Routes

- Payment processing endpoints (see `app/src/routers/stripe.go`)

### Admin Routes

- Administrative endpoints (see `app/src/routers/admin.go`)

## Lambda Email Service

The project includes an AWS Lambda function that processes payment confirmation emails via SQS. See `lambda/README_TERRAFORM.md` for detailed deployment instructions.

### Quick Lambda Deployment

```bash
cd lambda
export SMTP_PASSWORD=your_password
export AWS_REGION=us-east-1
make deploy-terraform
```

## Development

### Database Migrations

The application automatically runs migrations on startup for:
- `Play` model
- `Presentation` model
- `User` model

### Adding Dependencies

```bash
cd app
go get <package-name>
go mod tidy
```

## Configuration

### Database

The application uses PostgreSQL with GORM. Database connection is configured through environment variables in `.env`.

### Stripe

Configure your Stripe keys in `.env`:
- `STRIPE_PUBLIC_KEY`: Your Stripe publishable key
- `STRIPE_PRIVATE_KEY`: Your Stripe secret key

### AWS Lambda

For the email service, configure:
- `AWS_REGION`: AWS region for Lambda and SQS
- `SQS_QUEUE_URL`: URL of the SQS queue that triggers the Lambda
- `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`: AWS credentials

## Contributing

1. Create a feature branch
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## License

[Add your license here]

## Notes

- Redis is currently commented out in `docker-compose.yml` but can be enabled if needed
- The application includes a TODO for creating a thread to delete successful Redis tasks daily at 12:00 AM
- The Lambda function requires SMTP configuration (host, port, user, password) for sending emails

