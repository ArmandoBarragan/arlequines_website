package handlers

import (
	"math"
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/paymentintent"
)

type PaymentHandler interface {
	StripeWebhook(c *fiber.Ctx) error
	Success(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
}

type paymentHandler struct {
	paymentService      services.PaymentService
	presentationService services.PresentationService
}

func CreatePaymentHandler(
	paymentService services.PaymentService,
	presentationService services.PresentationService,
) paymentHandler {
	return paymentHandler{
		paymentService:      paymentService,
		presentationService: presentationService,
	}
}

func (handler paymentHandler) StripeWebhook(c *fiber.Ctx) error {
	/*
		Webhook to create a checkout session for a presentation
		Receives the amount of tickets, the presentation id and the email
		Creates a checkout session for the presentation
		Returns the url to the checkout session
	*/
	var body services.StripeWebhook

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	checkoutSession, err := handler.paymentService.CreateCheckoutSession(body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating checkout session"})
	}

	return c.Status(200).JSON(fiber.Map{"url": checkoutSession})
}

func (handler *paymentHandler) Success(c *fiber.Ctx) error {
	/* Fetch the information from the payment and store it in the database, as well as the
	reservation, then generate the SQS message that is sent to the lambda */
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Session ID was invalid"})
	}
	stripeSession, err := session.Get(sessionID, nil)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"error": "Error getting stripe session"})
	}
	paymentIntent, err := paymentintent.Get(stripeSession.PaymentIntent.ID, nil)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"error": "Error getting the payment intent"})
	}
	presentationID, err := strconv.Atoi(c.Query("presentation_id"))
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"error": "Error getting presentation"})
	}
	paidAmount := float64(paymentIntent.Amount) / 100.0
	presentation, err := handler.presentationService.GetPresentation(uint(presentationID))
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"error": err.Error()})
	}
	// Validate that it's actually a round number. Can't have reservations for 2.4 people
	var quantity float64 = paidAmount / presentation.Price
	if quantity != math.Trunc(quantity) {
		return c.Status(503).JSON(
			fiber.Map{"error": "Somehow the quantity is not a whole number. Please contact tech support"},
		)
	}
	payment := &models.Payment{
		SessionID:      sessionID,
		PresentationID: uint(presentationID),
		Email:          stripeSession.CustomerEmail,
		Amount:         paidAmount,
		Quantity:       int(quantity),
	}
	if err := handler.paymentService.CreateEmailSendingEventToSQS(payment); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	err = handler.paymentService.CreatePayment(payment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating payment"})
	}
	return c.SendStatus(200)
}

func (handler *paymentHandler) Cancel(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
