package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// PaymentEvent matches the struct from your main application
type PaymentEvent struct {
	Email          string `json:"email"`
	Amount         int64  `json:"amount"`
	Quantity       int64  `json:"quantity"`
	PresentationID uint   `json:"presentation_id"`
}

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	FromEmail string
}

// Handler processes SQS messages containing PaymentEvent
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	config := loadEmailConfig()

	for _, message := range sqsEvent.Records {
		var paymentEvent PaymentEvent

		// Parse the message body
		if err := json.Unmarshal([]byte(message.Body), &paymentEvent); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Send confirmation email
		if err := sendConfirmationEmail(paymentEvent, config); err != nil {
			log.Printf("Error sending email for presentation %d: %v", paymentEvent.PresentationID, err)
			return err
		}

		log.Printf("Successfully sent confirmation email for presentation %d to %s",
			paymentEvent.PresentationID, paymentEvent.Email)
	}
	return nil
}

// loadEmailConfig loads SMTP configuration from environment variables
func loadEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:  os.Getenv("SMTPHost"),
		SMTPPort:  os.Getenv("SMTPPort"),
		SMTPUser:  os.Getenv("SMTPUser"),
		SMTPPass:  os.Getenv("SMTPPass"),
		FromEmail: os.Getenv("FromEmail"),
	}
}

func sendConfirmationEmail(event PaymentEvent, config EmailConfig) error {
	amount := float64(event.Amount) / 100.0
	fmt.Println(amount)
	subject := fmt.Sprintf("Confirmación de Compra - Presentación %d", event.PresentationID)
	body := fmt.Sprintf("Hola!")
	message := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\n\r\n%s",
		event.Email,
		subject,
		body,
	)

	fmt.Println(
		"user: ", config.SMTPUser,
		"pass: ", config.SMTPPass,
		"host: ", config.SMTPHost,
		"port: ", config.SMTPPort,
		"from: ", config.FromEmail,
		"to: ", event.Email,
	)
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, config.SMTPHost)
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{event.Email}, []byte(message))

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
