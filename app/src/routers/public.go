package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupPublicRoutes(app *fiber.App, db *gorm.DB, secretKey string) {
	// Initialize repositories
	playRepository := repositories.NewPlayRepository(db)
	presentationRepository := repositories.NewPresentationRepository(db)

	// Initialize services
	playService := services.NewPlayService(playRepository)
	presentationService := services.NewPresentationService(presentationRepository)

	// Initialize handlers
	playHandler := handlers.CreatePlayHandler(playService)
	presentationHandler := handlers.CreatePresentationHandler(presentationService)

	// Play routes (public)
	app.Get("/plays", playHandler.GetList) // sorted alphabetically
	app.Get("/plays/:id", playHandler.GetDetail)

	// Presentation routes (public)
	app.Get("/presentations", presentationHandler.GetList) // sorted by date
	app.Get("/presentations/:id", presentationHandler.GetDetail)
}
