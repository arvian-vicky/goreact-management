package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
)

type CommentRepository interface {
	FindByCardID(cardID string) ([]models.Comment, error)
	FindByID(id uint) (*models.Comment, error)
	FindByPublicID(id string) (*models.Comment, error)
	Create(comment *models.Comment) error
	Update(comment *models.Comment) error
	Delete(id uint) error
}

type commentRepository struct{}

func NewCommentRepository() CommentRepository {
	return &commentRepository{}
}

func (r *commentRepository) FindByCardID(cardPublicID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Joins("JOIN cards ON cards.internal_id = comments.card_internal_id").
		Where("cards.public_id = ?", cardPublicID).
		Order("created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.First(&comment, id).Error
	return &comment, err
}

func (r *commentRepository) FindByPublicID(id string) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.Where("public_id = ?", id).First(&comment).Error
	return &comment, err
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return config.DB.Create(comment).Error
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return config.DB.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}
