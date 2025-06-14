package handlers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB // This should be set from main or via dependency injection

func SetDB(database *gorm.DB) {
	db = database
}

func ListPlays(c *fiber.Ctx) error {
	/*
		Lists all plays
	*/
	var plays []models.Play
	if err := db.Order("name asc").Find(&plays).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(plays)
}

func GetPlay(c *fiber.Ctx) error {
	/*
		Gets a play by id
	*/
	id := c.Params("id")
	var play models.Play
	if err := db.First(&play, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Play not found"})
	}
	return c.JSON(play)
}

func CreatePlay(c *fiber.Ctx) error {
	/*
		Creates a new play
		Receives the play
		Creates the play
	*/
	var play structs.Play
	if err := c.BodyParser(&play); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	model := models.Play{Name: play.Name, Author: play.Author}
	if err := db.Create(&model).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(model)
}

func UpdatePlay(c *fiber.Ctx) error {
	/*
		Updates a play
		Receives the id of the play and the play
		Updates the play
	*/
	id := c.Params("id")
	var play models.Play
	if err := db.First(&play, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Play not found"})
	}
	var input structs.Play
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	play.Name = input.Name
	play.Author = input.Author
	if err := db.Save(&play).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(play)
}

func DeletePlay(c *fiber.Ctx) error {
	/*
		Deletes a play
		Receives the id of the play
	*/
	id := c.Params("id")
	if err := db.Delete(&models.Play{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
