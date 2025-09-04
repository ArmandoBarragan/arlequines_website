package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAdminRoutes(app *fiber.App, db *gorm.DB, secretKey string) {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(db, secretKey)
	playRepository := repositories.NewPlayRepository(db)
	presentationRepository := repositories.NewPresentationRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepository, secretKey)
	playService := services.NewPlayService(playRepository)
	presentationService := services.NewPresentationService(presentationRepository)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	playHandler := handlers.CreatePlayHandler(playService)
	presentationHandler := handlers.CreatePresentationHandler(presentationService)

	admin := app.Group("/admin")
	admin.Use(authHandler.Protected())
	admin.Use(authHandler.AdminOnly())

	// Admin dashboard
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to admin dashboard",
		})
	})

	// Play management routes
	plays := admin.Group("/plays")
	plays.Post("/", playHandler.Create)
	plays.Put("/:id", playHandler.Update)
	plays.Delete("/:id", playHandler.Delete)

	// Presentation management routes
	presentations := admin.Group("/presentations")
	presentations.Post("/", presentationHandler.Create)
	presentations.Put("/:id", presentationHandler.Update)
	presentations.Delete("/:id", presentationHandler.Delete)
}
