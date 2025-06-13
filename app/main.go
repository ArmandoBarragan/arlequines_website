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

	db.AutoMigrate(&models.Play{}, &models.Presentation{}, &models.User{})

	handlers.SetDB(db)

	app := fiber.New()

	routers.SetupRoutes(app)
	routers.SetupAdminRoutes(app)
	routers.SetupAuthRoutes(app)

	app.Listen(":8000")
}
