package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

var Config *settings.Config

type PaymentService interface {
	CreateCheckoutSession(webhook StripeWebhook) (*stripe.CheckoutSession, error)
	CreatePayment(payment *models.Payment) error
	CreateEmailSendingEventToSQS(payment *models.Payment) error
}

type paymentService struct {
	presentationRepository repositories.PresentationRepository
	playRepository repositories.PlayRepository
	paymentRepository repositories.PaymentRepository
}

func NewPaymentService(
	presentationRepository repositories.PresentationRepository,
	playRepository repositories.PlayRepository,
	paymentRepository repositories.PaymentRepository,
) PaymentService {
	return &paymentService{
		presentationRepository: presentationRepository,
		playRepository: playRepository,
		paymentRepository: paymentRepository,
	}
}


func (service paymentService) CreatePayment(payment *models.Payment) error {
	return service.paymentRepository.Create(payment)
}

func (service paymentService) CreateEmailSendingEventToSQS(payment *models.Payment) error {
	localConfig := settings.LoadConfig()
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(localConfig.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			localConfig.AWSAccessKeyID,
			localConfig.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create SQS client
	sqsClient := sqs.NewFromConfig(awsConfig)
	// Convert PaymentEvent to JSON
	messageBody, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment event: %w", err)
	}

	// Send message to SQS
	_, err = sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(localConfig.SQSQueueURL),
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"EventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("PaymentSuccess"),
			},
			"PresentationName": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(payment.PresentationName),
			},
		},
		MessageGroupId: aws.String("payment_group"),
		MessageDeduplicationId: aws.String(fmt.Sprintf("%d", payment.ID)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	fmt.Printf("Successfully sent payment event to SQS for presentation %s", payment.PresentationName)
	return nil
}

type StripeWebhook struct {
	AmountOfTickets int    `json:"amount_of_tickets"`
	PresentationID  uint   `json:"presentation_id"`
	Email           string `json:"email"`
}

func (service paymentService) CreateCheckoutSession(webhook StripeWebhook) (*stripe.CheckoutSession, error) {
	presentation, err := service.presentationRepository.FindByID(uint(webhook.PresentationID))
	if err != nil {
		return nil, err
	}
	checkoutSession, err := service.createCheckoutSession(presentation, webhook)
	if err != nil {
		return nil, err
	}
	return checkoutSession, nil
}

func (service paymentService) createCheckoutSession(presentation *models.Presentation, webhook StripeWebhook) (*stripe.CheckoutSession, error) {
	config := settings.LoadConfig()
	successURL := "/stripe/success?session_id={CHECKOUT_SESSION_ID}&presentation_id="
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyMXN)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Ticket"), // Product name
					},
					UnitAmount: stripe.Int64(int64(presentation.Price * 100)), // Price in cents (2000 cents = $20.00)
				},
				Quantity: stripe.Int64(int64(webhook.AmountOfTickets)), // Quantity of the item
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)), // Set mode to 'payment' for one-time payments
		SuccessURL: stripe.String(
			config.HostURL + successURL + strconv.Itoa(int(presentation.ID)),
		),
		CancelURL:     stripe.String(config.HostURL + "/stripe/cancel"),
		CustomerEmail: stripe.String(webhook.Email),
	}
	return session.New(params)
}

