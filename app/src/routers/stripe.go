package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupStripeRoutes(app *fiber.App) {
	app.Post("/stripe/webhook", handlers.StripeWebhook)
	app.Get("/stripe/success", handlers.Success)
	app.Get("/stripe/cancel", handlers.Cancel)
}
