package handlers

import (
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler interface {
	StripeWebhook(c *fiber.Ctx) error
	Success(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
}

type paymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) paymentHandler {
	return paymentHandler{paymentService: paymentService}
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
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Session ID was invalid"})
	}
	presentationID, err := strconv.Atoi(c.Query("presentation_id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting presentation"})
	}
	amount, err := strconv.Atoi(c.Query("amount"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting amount"})
	}
	quantity, err := strconv.Atoi(c.Query("quantity"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting quantity"})
	}
	payment := &models.Payment{
		SessionID: sessionID,
		PresentationID: uint(presentationID),
		Email: c.Query("email"),
		Amount: float64(amount) / 100.0,
		Quantity: quantity,
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
