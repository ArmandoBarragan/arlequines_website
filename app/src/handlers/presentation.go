package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PresentationHandler interface {
	Create(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	GetList(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type presentationHandler struct {
	service services.PresentationService
}

func CreatePresentationHandler(service services.PresentationService) presentationHandler {
	return presentationHandler{service: service}
}

func (handler *presentationHandler) Create(c *fiber.Ctx) error {
	var presentation services.CreatePresentationsRequest
	if err := c.BodyParser(&presentation); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	presentationFromDB, err := handler.service.CreatePresentation(presentation)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(presentationFromDB)
}

func (handler *presentationHandler) parsePresentation(presentation *models.Presentation) services.PresentationResponse {
	return services.PresentationResponse{
		ID:             presentation.ID,
		PlayID:         presentation.PlayID,
		DateTime:       presentation.DateTime,
		Location:       presentation.Location,
		Price:          presentation.Price,
		SeatLimit:      presentation.SeatLimit,
		AvailableSeats: presentation.AvailableSeats,
	}
}

func (handler *presentationHandler) GetList(c *fiber.Ctx) error {
	filter := repositories.PresentationFilter{
		Limit:  10, // Default limit
		Offset: 0,  // Default offset
	}

	// Parse limit and offset
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

	// Parse date filter
	if dateStr := c.Query("date"); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			filter.Date = date
		}
	}

	// Parse location filter
	if location := c.Query("location"); location != "" {
		filter.Location = location
	}

	// Parse price filter
	if priceStr := c.Query("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil && price > 0 {
			filter.Price = price
		}
	}

	presentations, err := handler.service.ListPresentations(filter)
	var presentationResponse []services.PresentationResponse
	for _, presentation := range presentations {
		presentationResponse = append(presentationResponse, handler.parsePresentation(&presentation))
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(presentationResponse)
}

func (handler *presentationHandler) GetDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	presentation, err := handler.service.GetPresentation(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(handler.parsePresentation(presentation))

}

func (handler *presentationHandler) Update(c *fiber.Ctx) error {
	var presentation services.UpdatePresentationRequest
	if err := c.BodyParser(&presentation); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	presentationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	presentation.ID = uint(presentationID)
	updatedPresentation, err := handler.service.UpdatePresentation(presentation)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(handler.parsePresentation(&updatedPresentation))
}

func (handler *presentationHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	err = handler.service.DeletePresentation(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Presentation not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(204).JSON(fiber.Map{"message": "presentation deleted successfuly"})
}
