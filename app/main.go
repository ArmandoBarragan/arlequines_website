package main

import (
	"log"

	"github.com/ArmandoBarragan/arlequines_website/settings"
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

	// TODO: Create a thread that deletes successful redis tasks every day at 12:00 AM
	app := fiber.New()

	routers.SetupPublicRoutes(app, db, config.SecretKey)
	routers.SetupStripeRoutes(app, db, config.SecretKey)
	routers.SetupAuthRoutes(app, db, config.SecretKey)
	routers.SetupAdminRoutes(app, db, config.SecretKey)
	app.Listen(":8000")
}
