package handlers

import (
	"errors"
	"strconv"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PlayHandler interface {
	Create(c *fiber.Ctx) error
	GetList(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type playHandler struct {
	service services.PlayService
}

func CreatePlayHandler(service services.PlayService) *playHandler {
	return &playHandler{service: service}
}

func parsePlay(play *models.Play) services.PlayResponse {
	return services.PlayResponse{
		ID:     play.ID,
		Name:   play.Name,
		Author: play.Author,
	}
}

func (handler *playHandler) Create(c *fiber.Ctx) error {
	var play services.CreatePlayRequest
	if err := c.BodyParser(&play); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	playFromDB, err := handler.service.CreatePlay(&play)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(playFromDB)
}

func (handler *playHandler) GetList(c *fiber.Ctx) error {
	filter := repositories.PlayFilter{
		Limit:  10,
		Offset: 0,
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}
	if name := c.Query("name"); name != "" {
		filter.Name = name
	}
	if author := c.Query("author"); author != "" {
		filter.Name = author
	}
	plays, err := handler.service.ListPlays(&filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(plays)
}

func (handler *playHandler) GetDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	play, err := handler.service.GetPlay(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Play not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(play)

}

func (handler *playHandler) Update(c *fiber.Ctx) error {
	var play services.UpdatePlayRequest
	if err := c.BodyParser(&play); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	playID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	play.ID = uint(playID)
	updatedPlay, err := handler.service.UpdatePlay(&play)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Play not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(parsePlay(updatedPlay))
}

func (handler *playHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	err = handler.service.DeletePlay(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Play not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(204).JSON(fiber.Map{"message": "play deleted successfuly"})
}
