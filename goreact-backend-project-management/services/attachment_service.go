package services

import (
	"errors"
	"time"

	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"

	"github.com/google/uuid"
)

type AttachmentService interface {
	GetByCardID(cardPublicID string) ([]models.CardAttachment, error)
	GetByID(id uint) (*models.CardAttachment, error)
	Create(cardPublicID string, userPublicID string, filename string) (*models.CardAttachment, error)
	Delete(id uint) error
	GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error)
	DeleteByPublicID(pubID uuid.UUID) error
}

type attachmentService struct {
	attachmentRepo repositories.AttachmentRepository
	cardRepo       repositories.CardRepository
	userRepo       repositories.UserRepository
}

func NewAttachmentService(
	attachmentRepo repositories.AttachmentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
) AttachmentService {
	return &attachmentService{attachmentRepo, cardRepo, userRepo}
}

func (s *attachmentService) GetByCardID(cardPublicID string) ([]models.CardAttachment, error) {
	// Verifikasi card ada lewat cardRepo
	_, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}
	return s.attachmentRepo.FindByCardID(cardPublicID)
}

func (s *attachmentService) GetByID(id uint) (*models.CardAttachment, error) {
	return s.attachmentRepo.FindByID(id)
}

func (s *attachmentService) Create(cardPublicID string, userPublicID string, filename string) (*models.CardAttachment, error) {
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}
	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	attach := &models.CardAttachment{
		PublicID:  uuid.New(),
		CardID:    card.InternalID,
		UserID:    user.InternalID,
		File:      filename,
		CreatedAt: time.Now(),
	}

	if err := s.attachmentRepo.Create(attach); err != nil {
		return nil, err
	}
	return attach, nil
}

func (s *attachmentService) Delete(id uint) error {
	return s.attachmentRepo.Delete(id)
}

func (s *attachmentService) GetByPublicID(pubID uuid.UUID) (*models.CardAttachment, error) {
	return s.attachmentRepo.FindByPublicID(pubID)
}

func (s *attachmentService) DeleteByPublicID(pubID uuid.UUID) error {
	return s.attachmentRepo.DeleteByPublicID(pubID)
}
