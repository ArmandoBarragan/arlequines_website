package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", handlers.CreateAccount)
	auth.Post("/login", handlers.Login)

	// Protected user routes
	auth.Get("/me", handlers.Protected(), func(c *fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(user)
	})
}
