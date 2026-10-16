package services

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"
)

type LabelService interface {
	GetAll() ([]models.Label, error)
	GetByID(id uint) (*models.Label, error)
	GetByPublicID(publicID string) (*models.Label, error)
	Create(label *models.Label) error
	Update(label *models.Label) error
	Delete(id uint) error
}

type labelService struct {
	repo repositories.LabelRepository
}

func NewLabelService(repo repositories.LabelRepository) LabelService {
	return &labelService{repo}
}

func (s *labelService) GetAll() ([]models.Label, error) {
	return s.repo.FindAll()
}

func (s *labelService) GetByID(id uint) (*models.Label, error) {
	return s.repo.FindByID(id)
}

func (s *labelService) GetByPublicID(publicID string) (*models.Label, error) {
	return s.repo.FindByPublicID(publicID)
}

func (s *labelService) Create(label *models.Label) error {
	return s.repo.Create(label)
}

func (s *labelService) Update(label *models.Label) error {
	return s.repo.Update(label)
}

func (s *labelService) Delete(id uint) error {
	return s.repo.Delete(id)
}
