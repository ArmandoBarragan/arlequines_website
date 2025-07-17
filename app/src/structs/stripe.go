package structs

import (
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type StripeWebhook struct {
	AmountOfTickets int    `json:"amount_of_tickets"`
	PresentationID  int    `json:"presentation_id"`
	Email           string `json:"email"`
}

func (s *StripeWebhook) CreateCheckoutSession(presentation models.Presentation) (*stripe.CheckoutSession, error) {
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
				Quantity: stripe.Int64(int64(s.AmountOfTickets)), // Quantity of the item
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)), // Set mode to 'payment' for one-time payments
		SuccessURL: stripe.String(
			config.HostURL + successURL + strconv.Itoa(int(presentation.ID)),
		), // URL to redirect after successful payment
		CancelURL:     stripe.String(config.HostURL + "/stripe/cancel"), // URL to redirect if payment is cancelled
		CustomerEmail: stripe.String(s.Email),
	}
	return session.New(params)
}
