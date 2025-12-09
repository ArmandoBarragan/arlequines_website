package settings

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82"
)

type Config struct {
	SecretKey                   string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBPassword                  string
	DBName                      string
	DBSSLMode                   string
	HostURL                     string
	StripePublicKey             string
	StripePrivateKey            string
	SMTPHost                    string
	SMTPPort                    string
	SMTPUser                    string
	SMTPPass                    string
	ConsumerTimer               time.Duration
	AWSRegion          string
	SQSQueueURL        string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
}

func LoadConfig() *Config {
	timer, err := strconv.Atoi(os.Getenv("CONSUMER_TIMER"))
	if err != nil {
		log.Printf("Failed to load consumer timer")
		timer = 60
	}
	config := &Config{
		SecretKey:        os.Getenv("SECRET_KEY"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMode:        os.Getenv("DB_SSLMODE"),
		HostURL:          os.Getenv("HOST_URL"),
		StripePublicKey:  os.Getenv("STRIPE_PUBLIC_KEY"),
		StripePrivateKey: os.Getenv("STRIPE_PRIVATE_KEY"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPass:      os.Getenv("SMTP_PASS"),
		ConsumerTimer: time.Duration(timer),
		AWSRegion:          os.Getenv("AWS_REGION"),
		SQSQueueURL:        os.Getenv("SQS_QUEUE_URL"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	stripe.Key = config.StripePrivateKey
	return config
}
