package services

import (
	"fmt"
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type StripeService interface {
	CreateCheckoutSession(webhook StripeWebhook) (*stripe.CheckoutSession, error)
	ProcessPaymentSuccessRedis(presentationID uint, sessionID string) error
	ProcessPaymentSuccessSQS(presentationID uint, sessionID string) error
}

type stripeService struct {
	repository repositories.PresentationRepository
}

func NewStripeService(repository repositories.PresentationRepository) stripeService {
	return stripeService{repository: repository}
}

type StripeWebhook struct {
	AmountOfTickets int    `json:"amount_of_tickets"`
	PresentationID  uint   `json:"presentation_id"`
	Email           string `json:"email"`
}

func (service stripeService) CreateCheckoutSession(webhook StripeWebhook) (*stripe.CheckoutSession, error) {
	presentation, err := service.repository.FindByID(uint(webhook.PresentationID))
	if err != nil {
		return nil, err
	}
	checkoutSession, err := service.createCheckoutSession(presentation, webhook)
	if err != nil {
		return nil, err
	}
	return checkoutSession, nil
}

func (service stripeService) createCheckoutSession(presentation *models.Presentation, webhook StripeWebhook) (*stripe.CheckoutSession, error) {
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

func (service stripeService) ProcessPaymentSuccessRedis(presentationID uint, sessionID string) error {
	/* Decreases the available seats in the database, creates a QR and retrieves it
	trough HTTP, but also creates an email sending event that sends the same QR to the
	user's email*/
	presentation, err := service.repository.FindByID(presentationID)
	if err != nil {
		return err
	}
	sessionData, err := session.Get(sessionID, nil)
	if err != nil {
		return err
	}
	quantity := int64(float64(sessionData.AmountTotal) / presentation.Price)
	paymentEvent := PaymentEvent{
		Email:          sessionData.CustomerEmail,
		Amount:         sessionData.AmountTotal,
		Quantity:       quantity,
		PresentationID: presentationID,
	}
	paymentEvent.CreateEmailSendingEventRedis()
	return nil
}

func (service stripeService) ProcessPaymentSuccessSQS(presentationID uint, sessionID string) error {
	sessionData, err := session.Get(sessionID, nil)
	if err != nil {
		return err
	}
	presentation, err := service.repository.FindByID(presentationID)
	if err != nil {
		return err
	}
	quantity := int64(float64(sessionData.AmountTotal) / presentation.Price)
	paymentEvent := PaymentEvent{
		Email:          sessionData.CustomerEmail,
		Amount:         sessionData.AmountTotal,
		Quantity:       quantity,
		PresentationID: presentationID,
	}
	err = paymentEvent.CreateEmailSendingEventSQS()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully sent payment event to SQS for presentation %d", presentationID)
	return nil
}
