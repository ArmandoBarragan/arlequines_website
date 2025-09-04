package services

import (
	"github.com/ArmandoBarragan/arlequines_website/src/models"
	"github.com/ArmandoBarragan/arlequines_website/src/repositories"
)

type PlayResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type CreatePlayRequest struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type UpdatePlayRequest struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type PlayService interface {
	CreatePlay(play *CreatePlayRequest) (*models.Play, error)
	GetPlay(id uint) (*models.Play, error)
	ListPlays(filter *repositories.PlayFilter) ([]models.Play, error)
	UpdatePlay(play *UpdatePlayRequest) (*models.Play, error)
	DeletePlay(id uint) error
}

type playService struct {
	repository repositories.PlayRepository
}

func NewPlayService(repository repositories.PlayRepository) *playService {
	return &playService{repository: repository}
}

func (service *playService) CreatePlay(request *CreatePlayRequest) (*models.Play, error) {
	var play models.Play = models.Play{
		Name:   request.Name,
		Author: request.Author,
	}
	err := service.repository.Create(&play)
	if err != nil {
		return nil, err
	}
	return &play, nil
}

func (service *playService) GetPlay(id uint) (*models.Play, error) {
	play, err := service.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return play, nil
}

func (service *playService) ListPlays(playFilter *repositories.PlayFilter) ([]models.Play, error) {
	plays, err := service.repository.List(*playFilter)
	if err != nil {
		return nil, err
	}
	return plays, nil
}
func (service *playService) UpdatePlay(request *UpdatePlayRequest) (*models.Play, error) {
	existingPlay, err := service.repository.FindByID(request.ID)
	if err != nil {
		return nil, err
	}

	existingPlay.Name = request.Name
	existingPlay.Author = request.Author

	err = service.repository.Update(existingPlay)
	if err != nil {
		return nil, err
	}
	return existingPlay, nil
}
func (service *playService) DeletePlay(id uint) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
