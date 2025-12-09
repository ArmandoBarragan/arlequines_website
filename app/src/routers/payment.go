package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(app *fiber.App, db *gorm.DB, secretKey string) {
	// Initialize repositories and services as needed
	presentationRepository := repositories.NewPresentationRepository(db)
	playRepository := repositories.NewPlayRepository(db)
	paymentRepository := repositories.NewPaymentRepository(db)
	paymentService := services.NewPaymentService(presentationRepository, playRepository, paymentRepository)
	presentationService := services.NewPresentationService(presentationRepository)
	handler := handlers.CreatePaymentHandler(paymentService, presentationService)

	// Initialize handlers
	app.Post("/payment/webhook", handler.StripeWebhook)
	app.Get("/payment/success", handler.Success)
	app.Get("/payment/cancel", handler.Cancel)
}
