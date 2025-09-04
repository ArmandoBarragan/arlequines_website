package repositories

import (
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"gorm.io/gorm"
)

type PresentationRepository interface {
	Create(presentation models.Presentation) error
	CreateBatch(presentations []models.Presentation) error
	FindByID(id uint) (*models.Presentation, error)
	Update(presentation models.Presentation) error
	Delete(id uint) error
	List(filter PresentationFilter) ([]models.Presentation, error)
}

type PresentationFilter struct {
	Date     time.Time
	Location string
	Price    float64
	Limit    int
	Offset   int
}

type presentationRepository struct {
	db *gorm.DB
}

func NewPresentationRepository(db *gorm.DB) *presentationRepository {
	return &presentationRepository{db: db}
}

func (r *presentationRepository) List(f PresentationFilter) ([]models.Presentation, error) {
	var presentation []models.Presentation
	q := r.db.Model(&models.Presentation{})

	if !f.Date.IsZero() {
		q = q.Where("datetime >= ?", f.Date)
	}
	if f.Location != "" {
		q = q.Where("location ILIKE ?", "%"+f.Location+"%")
	}
	if f.Price > 0 {
		q = q.Where("price <= ?", f.Price)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	q = q.Offset(f.Offset)

	err := q.Find(&presentation).Error
	return presentation, err
}

func (r *presentationRepository) CreateBatch(presentations []models.Presentation) error {
	return r.db.Create(&presentations).Error
}

func (r *presentationRepository) Update(presentation models.Presentation) error {
	return r.db.Save(presentation).Error
}

func (r *presentationRepository) Create(presentation models.Presentation) error {
	return r.db.Create(presentation).Error
}

func (r *presentationRepository) FindByID(id uint) (*models.Presentation, error) {
	var presentation models.Presentation
	err := r.db.First(&presentation, id).Error
	if err != nil {
		return nil, err
	}
	return &presentation, nil
}

func (r *presentationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Presentation{}, id).Error
}
