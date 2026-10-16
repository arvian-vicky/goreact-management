package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
)

type BoardMemberRepository interface {
	GetMembers(boardID string) ([]models.User, error)
	IsMember(boardID, userID string) (bool, error)
	CountMembers(boardID string) (int64, error)
}

type boardMemberRepository struct{}

func NewBoardMemberRepository() BoardMemberRepository {
	return &boardMemberRepository{}
}

func (r *boardMemberRepository) GetMembers(boardPublicID string) ([]models.User, error) {
	var users []models.User
	err := config.DB.Joins("JOIN board_members ON board_members.user_internal_id = users.internal_id").
		Joins("JOIN boards ON boards.internal_id = board_members.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).
		Find(&users).Error
	return users, err
}

func (r *boardMemberRepository) IsMember(boardPublicID, userPublicID string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.BoardMember{}).
		Joins("JOIN boards ON boards.internal_id = board_members.board_internal_id").
		Joins("JOIN users ON users.internal_id = board_members.user_internal_id").
		Where("boards.public_id = ? AND users.public_id = ?", boardPublicID, userPublicID).
		Count(&count).Error
	return count > 0, err
}

func (r *boardMemberRepository) CountMembers(boardPublicID string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.BoardMember{}).
		Joins("JOIN boards ON boards.internal_id = board_members.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).
		Count(&count).Error
	return count, err
}
