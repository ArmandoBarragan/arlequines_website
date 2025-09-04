package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB, secretKey string) {
	userRepository := repositories.NewUserRepository(db, secretKey)
	authService := services.NewAuthService(userRepository, secretKey)
	authHandler := handlers.NewAuthHandler(authService)

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.CreateAccount)
	auth.Post("/login", authHandler.Login)

	// Protected user routes
	auth.Get("/me", authHandler.Protected(), func(c *fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(user)
	})
}
