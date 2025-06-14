package handlers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func StripeWebhook(c *fiber.Ctx) error {
	var body structs.StripeWebhook

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	presentation := models.Presentation{ID: uint(body.PresentationID)}
	if err := db.First(&presentation).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
	}

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
				Quantity: stripe.Int64(int64(body.AmountOfTickets)), // Quantity of the item
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)), // Set mode to 'payment' for one-time payments
		SuccessURL: stripe.String("http://localhost:4242/success"),           // URL to redirect after successful payment
		CancelURL:  stripe.String("http://localhost:4242/cancel"),            // URL to redirect if payment is cancelled
	}
	result, err := session.New(params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating checkout session"})
	}

	return c.Status(200).JSON(fiber.Map{"url": result.URL})
}
