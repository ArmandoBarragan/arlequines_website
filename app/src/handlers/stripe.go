package handlers

import (
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
)

type StripeHandler interface {
	StripeWebhook(c *fiber.Ctx) error
	Success(c *fiber.Ctx) error
}

type stripeHandler struct {
	service services.StripeService
}

func NewStripeHandler(service services.StripeService) stripeHandler {
	return stripeHandler{service: service}
}

func (handler stripeHandler) StripeWebhook(c *fiber.Ctx) error {
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

	checkoutSession, err := handler.service.CreateCheckoutSession(body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating checkout session"})
	}

	return c.Status(200).JSON(fiber.Map{"url": checkoutSession})
}

func (handler *stripeHandler) Success(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Session ID was invalid"})
	}
	presentationID, err := strconv.Atoi(c.Query("presentation_id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting presentation"})
	}

	err = handler.service.ProcessPaymentSuccess(uint(presentationID), sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(200)
}

func Cancel(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
