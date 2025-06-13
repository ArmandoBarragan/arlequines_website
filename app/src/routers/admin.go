package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {
	// Admin routes
	admin := app.Group("/admin")
	admin.Use(handlers.Protected())
	admin.Use(handlers.AdminOnly())

	// Admin dashboard
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to admin dashboard",
		})
	})

	// Play management routes
	plays := admin.Group("/plays")
	plays.Post("/", handlers.CreatePlay)
	plays.Put("/:id", handlers.UpdatePlay)
	plays.Delete("/:id", handlers.DeletePlay)

	// Presentation management routes
	presentations := admin.Group("/presentations")
	presentations.Post("/", handlers.CreatePresentations)
	presentations.Put("/:id", handlers.UpdatePresentation)
	presentations.Delete("/:id", handlers.DeletePresentation)
}
