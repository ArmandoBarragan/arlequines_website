package repositories

import (
	"fmt"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db        *gorm.DB
	secretKey string
}

type UserFilter struct {
	Name   string
	Email  string
	Limit  int
	Offset int
}

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(filter *UserFilter) ([]models.User, error)
	PasswordIsCorrect(password string, user *models.User) bool
}

func NewUserRepository(db *gorm.DB, secretKey string) *userRepository {
	return &userRepository{
		db:        db,
		secretKey: secretKey,
	}
}

func (r *userRepository) List(f *UserFilter) ([]models.User, error) {
	var user []models.User
	q := r.db.Model(&models.User{})
	if f.Name != "" {
		q.Where("name ilike ", f.Name)
	}
	if f.Email != "" {
		q.Where("email ilike", f.Email)
	}
	if f.Limit > 0 {
		q.Limit(f.Limit)
	}
	q.Offset(f.Offset)
	err := q.Find(&user).Error
	return user, err
}

func (userRepo *userRepository) Create(user *models.User) error {
	return userRepo.db.Create(user).Error
}

func (userRepo *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := userRepo.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := userRepo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) Update(user *models.User) error {
	return userRepo.db.Save(user).Error
}

func (userRepo *userRepository) Delete(id uint) error {
	return userRepo.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) PasswordIsCorrect(password string, user *models.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	fmt.Println(err)
	return err == nil
}
