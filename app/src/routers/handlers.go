package routers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
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

type PlayHandler interface {
	Create(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	GetList(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type playHandler struct {
	service services.PlayService
}

type AuthHandler interface {
	CreateAccount(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Protected() fiber.Handler
	AdminOnly() fiber.Handler
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *authHandler {
	return &authHandler{service: service}
}

func NewPresentationHandler(service services.PresentationService) *presentationHandler {
	return &presentationHandler{service: service}
}

func NewPlayHandler(service services.PlayService) *playHandler {
	return &playHandler{service: service}
}
