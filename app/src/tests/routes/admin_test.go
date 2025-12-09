package route_tests

import (
	"net/http/httptest"
	"testing"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/ArmandoBarragan/arlequines_website/src/routers"
	"github.com/ArmandoBarragan/arlequines_website/src/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func SetupAdminTest(t *testing.T) (*fiber.App, string) {
	// This function sets up the admin test environment
	app, db, secretKey := tests.SetupTestApp()
	routers.SetupAdminRoutes(app, db, secretKey)
	user := models.User{
		Email: "admin@example.com",
		Password: "password",
		Name: "Admin",
		IsAdmin: true,
	}
	db.Create(&user)
	userRepository := repositories.NewUserRepository(db, secretKey)
	authService := services.NewAuthService(userRepository, secretKey)
	response, err := authService.Login(services.LoginRequest{
		Email:    "admin@example.com",
		Password: "password",
	})
	if err != nil || response == nil || response.Token == "" {
		t.Fatalf("Failed to login: %v", err)
	}
	return app, response.Token
}

func TestAdminRoutes_Dashboard_WithAdminToken(t *testing.T) {
	// This test should return a successful status code because the user is authenticated
	app, token := SetupAdminTest(t)
	request := httptest.NewRequest("GET", "/admin/dashboard", nil)
	request.Header.Set("Authorization", "Bearer " + token)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestAdminRoutes_Dashboard_WithoutToken(t *testing.T) {
	// This test should return an unauthorized status code because the user is not authenticated
	app, _ := SetupAdminTest(t)
	request := httptest.NewRequest("GET", "/admin/dashboard", nil)
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
}
