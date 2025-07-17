package handlers

import (
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/structs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // In production, use environment variable

func CreateAccount(c *fiber.Ctx) error {
	/*
		Creates a new user
		Receives the user
		Creates the user
		Returns the created user
	*/
	var req structs.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists",
		})
	}

	// Create new user
	user := models.User{
		Email:    req.Email,
		Password: req.Password, // Will be hashed by BeforeSave hook
		Name:     req.Name,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(structs.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(c *fiber.Ctx) error {
	/*
		Logs in a user
		Receives the user
		Logs in the user
		Returns the logged in user
	*/
	var req structs.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user by email
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(structs.AuthResponse{
		Token: token,
		User:  user,
	})
}

func generateToken(user models.User) (string, error) {
	/*
		Generates a JWT token for a user
		Receives the user
		Generates a JWT token for the user
		Returns the JWT token
	*/
	// Create the Claims
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	return token.SignedString(jwtSecret)
}

// Middleware to protect routes
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		/*
			Protects a route
			Receives the request
			Protects the route
			Returns the protected route
		*/
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Parse token
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !parsedToken.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Add user info to context
		c.Locals("user", claims)
		return c.Next()
	}
}

// Middleware to check if user is admin
func AdminOnly() fiber.Handler {
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
