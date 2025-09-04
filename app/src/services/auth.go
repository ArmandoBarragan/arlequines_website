package services

import (
	"errors"
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type User struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsAdmin   bool      `json:"is_admin"`
}

type AuthService interface {
	CreateAccount(request CreateAccountRequest) (*AuthResponse, error)
	Login(request LoginRequest) (*AuthResponse, error)
	ParseAndValidateToken(tokenString string) (jwt.MapClaims, error)
}

type authService struct {
	repository repositories.UserRepository
	jwtSecret  []byte
}

func NewAuthService(repo repositories.UserRepository, secretKey string) AuthService {
	return &authService{
		repository: repo,
		jwtSecret:  []byte(secretKey),
	}
}

func (service authService) parseUser(user *models.User) *User {
	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Name:      user.Name,
		IsAdmin:   user.IsAdmin,
	}
}

func (service authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.jwtSecret)
}

func (service *authService) ParseAndValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return service.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (service *authService) CreateAccount(request CreateAccountRequest) (*AuthResponse, error) {
	existentUser, err := service.repository.FindByEmail(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existentUser != nil {
		return nil, errors.New("a user with that email already exists")
	}

	newUser := &models.User{
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
	}

	err = service.repository.Create(newUser)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := service.generateToken(newUser)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &AuthResponse{Token: token, User: service.parseUser(newUser)}, nil
}

func (service *authService) Login(request LoginRequest) (*AuthResponse, error) {
	user, err := service.repository.FindByEmail(request.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}

	// Check password
	if !service.repository.PasswordIsCorrect(request.Password, user) {
		return nil, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := service.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &AuthResponse{Token: token, User: service.parseUser(user)}, nil
}
