package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupStripeRoutes(app *fiber.App, db *gorm.DB, secretKey string) {
	// Initialize repositories and services as needed
	repository := repositories.NewPresentationRepository(db)
	service := services.NewStripeService(repository)
	handler := handlers.NewStripeHandler(service)

	// Initialize handlers
	app.Post("/stripe/webhook", handler.StripeWebhook)
	app.Post("/stripe/success", handler.Success)
	app.Get("/stripe/cancel", handlers.Cancel)
}
