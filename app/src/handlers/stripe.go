package handlers

import (
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func StripeWebhook(c *fiber.Ctx) error {
	/*
		Webhook to create a checkout session for a presentation
		Receives the amount of tickets, the presentation id and the email
		Creates a checkout session for the presentation
		Returns the url to the checkout session
	*/
	var body structs.StripeWebhook

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	presentation := models.Presentation{ID: uint(body.PresentationID)}
	if err := db.First(&presentation).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
	}

	checkoutSession, err := body.CreateCheckoutSession(presentation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating checkout session"})
	}

	return c.Status(200).JSON(fiber.Map{"url": checkoutSession})
}

func Success(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Session ID was invalid"})
	}
	presentationID, err := strconv.Atoi(c.Query("presentation_id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting presentation"})
	}
	presentation := models.Presentation{ID: uint(presentationID)}
	if err := db.First(&presentation).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
	}
	sessionData, err := session.Get(sessionID, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error getting session"})
	}
	quantity := int64(float64(sessionData.AmountTotal) / presentation.Price)
	paymentEvent := structs.PaymentEvent{
		Email:          sessionData.CustomerEmail,
		Amount:         sessionData.AmountTotal,
		Quantity:       quantity,
		PresentationID: presentationID,
	}
	paymentEvent.CreateEmailSendingEvent()
	return c.SendStatus(200)
}

func Cancel(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
