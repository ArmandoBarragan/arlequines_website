package repositories

import (
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id uint) (*models.Payment, error)
	FindByEmail(email string) ([]models.Payment, error)
	FindByPresentationID(presentationID uint) ([]models.Payment, error)
	Update(payment *models.Payment) error
	Delete(id uint) error
	List(filter PaymentFilter) ([]models.Payment, error)
}

type PaymentFilter struct {
	Email          string
	PresentationID uint
	MinAmount      float64
	MaxAmount      float64
	StartDate      time.Time
	EndDate        time.Time
	Limit          int
	Offset         int
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Presentation").First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByEmail(email string) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("email = ?", email).Preload("Presentation").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) FindByPresentationID(presentationID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("presentation_id = ?", presentationID).Preload("Presentation").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}

func (r *paymentRepository) List(f PaymentFilter) ([]models.Payment, error) {
	var payments []models.Payment
	q := r.db.Model(&models.Payment{}).Preload("Presentation")

	if f.Email != "" {
		q = q.Where("email ILIKE ?", "%"+f.Email+"%")
	}
	if f.PresentationID > 0 {
		q = q.Where("presentation_id = ?", f.PresentationID)
	}
	if f.MinAmount > 0 {
		q = q.Where("amount >= ?", f.MinAmount)
	}
	if f.MaxAmount > 0 {
		q = q.Where("amount <= ?", f.MaxAmount)
	}
	if !f.StartDate.IsZero() {
		q = q.Where("created_at >= ?", f.StartDate)
	}
	if !f.EndDate.IsZero() {
		q = q.Where("created_at <= ?", f.EndDate)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	q = q.Offset(f.Offset)

	err := q.Find(&payments).Error
	return payments, err
}

