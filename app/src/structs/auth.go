package structs

import "github.com/ArmandoBarragan/arlequines_website/src/models"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
