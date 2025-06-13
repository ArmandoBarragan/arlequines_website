package main

import (
	"log"

	"github.com/ArmandoBarragan/arlequines_website/settings"
	"github.com/ArmandoBarragan/arlequines_website/src/handlers"
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/routers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := settings.LoadConfig()
	db, err := settings.SetupDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Play{}, &models.Presentation{})

	handlers.SetDB(db)

	app := fiber.New()

	routers.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Listen(":8000")
}
