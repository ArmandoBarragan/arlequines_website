package repositories

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"gorm.io/gorm"
)

type PlayRepository interface {
	Create(play *models.Play) error
	FindByID(id uint) (*models.Play, error)
	Update(play *models.Play) error
	Delete(id uint) error
	List(filter PlayFilter) ([]models.Play, error)
}

type PlayFilter struct {
	Name   string
	Author string
	Limit  int
	Offset int
}

type playRepository struct {
	db *gorm.DB
}

func NewPlayRepository(db *gorm.DB) playRepository {
	return playRepository{db: db}
}

func (playRepo playRepository) Create(play *models.Play) error {
	err := playRepo.db.Create(play).Error
	return err
}

func (r playRepository) FindByID(id uint) (*models.Play, error) {
	var play *models.Play
	err := r.db.First(&play, id).Error
	if err != nil {
		return nil, err
	}
	return play, nil
}

func (r playRepository) Update(play *models.Play) error {
	return r.db.Save(play).Error
}

func (r playRepository) Delete(id uint) error {
	return r.db.Delete(&models.Play{}, id).Error
}

func (r playRepository) List(f PlayFilter) ([]models.Play, error) {
	var play []models.Play
	q := r.db.Model(&models.Play{})

	if f.Name != "" {
		q = q.Where("name ILIKE ?", "%"+f.Name+"%")
	}
	if f.Author != "" {
		q = q.Where("author ILIKE ?", "%"+f.Author+"%")
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	}
	q = q.Offset(f.Offset)

	err := q.Find(&play).Error
	return play, err
}
