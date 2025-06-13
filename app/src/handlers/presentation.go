package handlers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
)

// db is declared in play.go and should be set via SetDB

func ListPresentations(c *fiber.Ctx) error {
	var presentations []models.Presentation
	if err := db.Order("date_time asc").Find(&presentations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presentations)
}

func GetPresentation(c *fiber.Ctx) error {
	id := c.Params("id")
	var presentation models.Presentation
	if err := db.First(&presentation, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
	}
	return c.JSON(presentation)
}

func CreatePresentations(c *fiber.Ctx) error {
	var input structs.PresentationsList
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	var created []models.Presentation
	for _, p := range input.Presentations {
		model := models.Presentation{
			PlayID:         p.PlayID,
			DateTime:       p.DateTime,
			Location:       p.Location,
			Price:          p.Price,
			SeatLimit:      p.SeatLimit,
			AvailableSeats: p.AvailableSeats,
		}
		if err := db.Create(&model).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		created = append(created, model)
	}
	return c.Status(201).JSON(created)
}

func UpdatePresentation(c *fiber.Ctx) error {
	id := c.Params("id")
	var presentation models.Presentation
	if err := db.First(&presentation, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
	}
	var input structs.Presentation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	presentation.PlayID = input.PlayID
	presentation.DateTime = input.DateTime
	presentation.Location = input.Location
	presentation.Price = input.Price
	presentation.SeatLimit = input.SeatLimit
	presentation.AvailableSeats = input.AvailableSeats
	if err := db.Save(&presentation).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(presentation)
}

func DeletePresentation(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&models.Presentation{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
