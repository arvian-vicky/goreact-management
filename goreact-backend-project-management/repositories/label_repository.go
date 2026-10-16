package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
)

type LabelRepository interface {
	FindAll() ([]models.Label, error)
	FindByID(id uint) (*models.Label, error)
	FindByPublicID(publicID string) (*models.Label, error)
	Create(label *models.Label) error
	Update(label *models.Label) error
	Delete(id uint) error
}

type labelRepository struct{}

func NewLabelRepository() LabelRepository {
	return &labelRepository{}
}

func (r *labelRepository) FindAll() ([]models.Label, error) {
	var labels []models.Label
	err := config.DB.Find(&labels).Error
	return labels, err
}

func (r *labelRepository) FindByID(id uint) (*models.Label, error) {
	var label models.Label
	err := config.DB.First(&label, id).Error
	return &label, err
}

func (r *labelRepository) FindByPublicID(publicID string) (*models.Label, error) {
	var label models.Label
	err := config.DB.Where("public_id = ?", publicID).First(&label).Error
	return &label, err
}

func (r *labelRepository) Create(label *models.Label) error {
	return config.DB.Create(label).Error
}

func (r *labelRepository) Update(label *models.Label) error {
	return config.DB.Save(label).Error
}

func (r *labelRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Label{}, id).Error
}
