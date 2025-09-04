package services

import (
	"time"

	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
)

type PresentationResponse struct {
	ID             uint      `json:"id"`
	PlayID         uint      `json:"play_id"`
	DateTime       time.Time `json:"datetime"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	SeatLimit      int       `json:"seat_limit"`
	AvailableSeats int       `json:"available_seats"`
}

type PresentationsListResponse struct {
	Presentations []PresentationResponse `json:"presentations"`
}

type CreatePresentationRequest struct {
	PlayID         uint      `json:"play_id"`
	DateTime       time.Time `json:"datetime"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	SeatLimit      int       `json:"seat_limit"`
	AvailableSeats int       `json:"available_seats"`
}

type UpdatePresentationRequest struct {
	ID             uint      `json:"id"`
	PlayID         uint      `json:"play_id"`
	DateTime       time.Time `json:"datetime"`
	Location       string    `json:"location"`
	Price          float64   `json:"price"`
	SeatLimit      int       `json:"seat_limit"`
	AvailableSeats int       `json:"available_seats"`
}

type CreatePresentationsRequest struct {
	Presentations []CreatePresentationRequest `json:"presentations"`
}

type PresentationService interface {
	CreatePresentation(request CreatePresentationsRequest) ([]models.Presentation, error)
	GetPresentation(id uint) (*models.Presentation, error)
	ListPresentations(filter repositories.PresentationFilter) ([]models.Presentation, error)
	UpdatePresentation(presentation UpdatePresentationRequest) (models.Presentation, error)
	DeletePresentation(id uint) error
}

type presentationService struct {
	repository repositories.PresentationRepository
}

func NewPresentationService(repository repositories.PresentationRepository) *presentationService {
	return &presentationService{repository: repository}
}

func (service *presentationService) CreatePresentation(request CreatePresentationsRequest) ([]models.Presentation, error) {
	var presentations []models.Presentation

	for _, p := range request.Presentations {
		presentation := models.Presentation{
			PlayID:         p.PlayID,
			DateTime:       p.DateTime,
			Location:       p.Location,
			Price:          p.Price,
			SeatLimit:      p.SeatLimit,
			AvailableSeats: p.AvailableSeats,
		}
		presentations = append(presentations, presentation)
	}

	err := service.repository.CreateBatch(presentations)
	if err != nil {
		return nil, err
	}

	return presentations, nil
}

func (service *presentationService) GetPresentation(id uint) (*models.Presentation, error) {
	presentation, err := service.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return presentation, nil
}

func (service *presentationService) ListPresentations(filter repositories.PresentationFilter) ([]models.Presentation, error) {
	presentations, err := service.repository.List(filter)
	if err != nil {
		return nil, err
	}
	return presentations, nil
}

func (service *presentationService) UpdatePresentation(presentation UpdatePresentationRequest) (models.Presentation, error) {
	existingPresentation, err := service.repository.FindByID(presentation.ID)
	if err != nil {
		return models.Presentation{}, err
	}
	existingPresentation.PlayID = presentation.PlayID
	existingPresentation.DateTime = presentation.DateTime
	existingPresentation.Location = presentation.Location
	existingPresentation.Price = presentation.Price
	existingPresentation.SeatLimit = presentation.SeatLimit
	existingPresentation.AvailableSeats = presentation.AvailableSeats
	err = service.repository.Update(*existingPresentation)
	if err != nil {
		return models.Presentation{}, err
	}
	return *existingPresentation, nil
}

func (service *presentationService) DeletePresentation(id uint) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
