package route_tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ArmandoBarragan/arlequines_website/src/tests"
	"github.com/ArmandoBarragan/arlequines_website/src/routers"
	"github.com/ArmandoBarragan/arlequines_website/src/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func SetupAuthTest(t *testing.T) *fiber.App {
	// This function sets up the auth test environment
	app, db, secretKey := tests.SetupTestApp()
	routers.SetupAuthRoutes(app, db, secretKey)
	return app
}

func TestAuthRoutes_Register_Success(t *testing.T) {
	app := SetupAuthTest(t)
	reqBody := services.CreateAccountRequest{
		Email: "test@example.com",
		Password: "password",
		Name: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestAuthRoutes_Login_Success(t *testing.T) {
	app := SetupAuthTest(t)
	reqBody := services.CreateAccountRequest{
		Email: "test@example.com",
		Password: "password",
		Name: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	loginReqBody := services.LoginRequest{
		Email: "test@example.com",
		Password: "password",
	}
	loginJsonBody, _ := json.Marshal(loginReqBody)
	request = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(loginJsonBody))
	request.Header.Set("Content-Type", "application/json")
	response, err = app.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestAuthRoutes_Me_Success(t *testing.T) {
	// This test should return a successful status code because the user is authenticated
	app := SetupAuthTest(t)
	reqBody := services.CreateAccountRequest{
		Email: "test2@example.com",
		Password: "password",
		Name: "Test User 2",
	}
	jsonBody, _ := json.Marshal(reqBody)
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.NoError(t, err)

	// Extract token from response
	var authResponse services.AuthResponse
	json.NewDecoder(response.Body).Decode(&authResponse)
	token := authResponse.Token

	request = httptest.NewRequest("GET", "/auth/me", nil)
	request.Header.Set("Authorization", "Bearer " + token)
	response, err = app.Test(request)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}
