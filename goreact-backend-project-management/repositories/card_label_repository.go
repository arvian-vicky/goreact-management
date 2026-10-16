package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
)

type CardLabelRepository interface {
	GetLabels(cardID string) ([]models.Label, error)
	CardHasLabel(cardID, labelID string) (bool, error)
}

type cardLabelRepository struct{}

func NewCardLabelRepository() CardLabelRepository {
	return &cardLabelRepository{}
}

func (r *cardLabelRepository) GetLabels(cardPublicID string) ([]models.Label, error) {
	var labels []models.Label
	err := config.DB.Joins("JOIN card_labels ON card_labels.label_internal_id = labels.internal_id").
		Joins("JOIN cards ON cards.internal_id = card_labels.card_internal_id").
		Where("cards.public_id = ?", cardPublicID).
		Find(&labels).Error
	return labels, err
}

func (r *cardLabelRepository) CardHasLabel(cardPublicID, labelPublicID string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.CardLabel{}).
		Joins("JOIN cards ON cards.internal_id = card_labels.card_internal_id").
		Joins("JOIN labels ON labels.internal_id = card_labels.label_internal_id").
		Where("cards.public_id = ? AND labels.public_id = ?", cardPublicID, labelPublicID).
		Count(&count).Error
	return count > 0, err
}
