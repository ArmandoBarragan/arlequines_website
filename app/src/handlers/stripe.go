package handlers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
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
	return c.SendStatus(200)
}

func Cancel(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
