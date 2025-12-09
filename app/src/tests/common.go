package tests

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func SetupTestApp() (*fiber.App, *gorm.DB, string) {
	app := fiber.New()
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Play{})
	db.AutoMigrate(&models.Presentation{})
	db.AutoMigrate(&models.Payment{})
	secretKey := "secret"
	return app, db, secretKey
}
