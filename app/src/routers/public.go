package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Play routes (public)
	app.Get("/plays", handlers.ListPlays) // sorted alphabetically
	app.Get("/plays/:id", handlers.GetPlay)

	// Presentation routes (public)
	app.Get("/presentations", handlers.ListPresentations) // sorted by date
	app.Get("/presentations/:id", handlers.GetPresentation)

}
