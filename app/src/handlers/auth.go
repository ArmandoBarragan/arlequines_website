package handlers

import (
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler interface {
	CreateAccount(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Protected(c *fiber.Ctx) fiber.Handler
	AdminOnly(c *fiber.Ctx) fiber.Handler
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) authHandler {
	return authHandler{service: service}
}

func (handler *authHandler) CreateAccount(c *fiber.Ctx) error {
	var req services.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	response, err := handler.service.CreateAccount(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	response, err := handler.service.Login(req)
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// Middleware to protect routes
func (handler *authHandler) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := handler.service.ParseAndValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("user", claims)
		return c.Next()
	}
}

// Middleware to check if user is admin
func (handler *authHandler) AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		/*
			Checks if the user is an admin
			Receives the request
			Checks if the user is an admin
			Returns the admin only route
		*/
		user := c.Locals("user").(jwt.MapClaims)
		if !user["is_admin"].(bool) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}
