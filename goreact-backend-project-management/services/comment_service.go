package services

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"
	"errors"
)

type CommentService interface {
	GetByCardID(cardPublicID string) ([]models.Comment, error)
	GetByID(id uint) (*models.Comment, error)
	GetByPublicID(id string) (*models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	cardRepo    repositories.CardRepository
	userRepo    repositories.UserRepository
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	cardRepo repositories.CardRepository,
	userRepo repositories.UserRepository,
) CommentService {
	return &commentService{commentRepo, cardRepo, userRepo}
}

func (s *commentService) GetByCardID(cardPublicID string) ([]models.Comment, error) {
	// Verifikasi card ada
	_, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	return s.commentRepo.FindByCardID(cardPublicID)
}

func (s *commentService) GetByID(id uint) (*models.Comment, error) {
	return s.commentRepo.FindByID(id)
}

func (s *commentService) GetByPublicID(id string) (*models.Comment, error) {
	return s.commentRepo.FindByPublicID(id)
}

func (s *commentService) Create(comment *models.Comment) error {
	// Pastikan card dan user ada
	_, err := s.cardRepo.FindCardByPublicID(comment.CardPubID.String())
	if err != nil {
		return errors.New("card not found")
	}

	_, err = s.userRepo.FindByPublicID(comment.UserPubID.String())
	if err != nil {
		return errors.New("user not found")
	}

	return s.commentRepo.Create(comment)
}

func (s *commentService) Update(comment *models.Comment) error {
	return s.commentRepo.Update(comment)
}

func (s *commentService) Delete(id uint) error {
	return s.commentRepo.Delete(id)
}
