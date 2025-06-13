package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Play routes
	app.Get("/plays", handlers.ListPlays) // sorted alphabetically
	app.Get("/plays/:id", handlers.GetPlay)
	app.Post("/plays", handlers.CreatePlay)
	app.Put("/plays/:id", handlers.UpdatePlay)
	app.Delete("/plays/:id", handlers.DeletePlay)

	// Presentation routes
	app.Get("/presentations", handlers.ListPresentations) // sorted by date
	app.Get("/presentations/:id", handlers.GetPresentation)
	app.Post("/presentations", handlers.CreatePresentations) // accepts PresentationsList
	app.Put("/presentations/:id", handlers.UpdatePresentation)
	app.Delete("/presentations/:id", handlers.DeletePresentation)
}
