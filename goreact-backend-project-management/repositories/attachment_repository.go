package repositories

import (
	"errors"

	"github.com/google/uuid"

	"arfian-nurdavicky/project-management/models"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	FindByCardID(cardPublicID string) ([]models.CardAttachment, error)
	GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error)
	FindByID(id uint) (*models.CardAttachment, error)
	Create(attachment *models.CardAttachment) error
	Delete(id uint) error
	DeleteByPublicID(pubID uuid.UUID) error
	FindByPublicID(pubID uuid.UUID) (*models.CardAttachment, error)
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db}
}

func (r *attachmentRepository) FindByCardID(cardPublicID string) ([]models.CardAttachment, error) {
	// Ambil internal_id card dulu
	var card models.Card
	if err := r.db.Where("public_id = ?", cardPublicID).First(&card).Error; err != nil {
		return nil, err
	}

	var attachments []models.CardAttachment
	if err := r.db.Where("card_internal_id = ?", card.InternalID).Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *attachmentRepository) GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	var att models.CardAttachment
	if err := r.db.Where("public_id = ?", pubID).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}

func (r *attachmentRepository) FindByID(id uint) (*models.CardAttachment, error) {
	var a models.CardAttachment
	if err := r.db.Where("internal_id = ?", id).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *attachmentRepository) Create(attachment *models.CardAttachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) Delete(id uint) error {
	return r.db.Where("internal_id = ?", id).Delete(&models.CardAttachment{}).Error
}

func (r *attachmentRepository) FindByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	var att models.CardAttachment
	if err := r.db.Where("public_id = ?", pubID).First(&att).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}

func (r *attachmentRepository) DeleteByPublicID(pubID uuid.UUID) error {
	return r.db.Where("public_id = ?", pubID).Delete(&models.CardAttachment{}).Error
}
