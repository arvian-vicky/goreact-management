package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"

	"github.com/google/uuid"
)

type ListRepository interface {
	FindByBoardID(boardID string) ([]models.List, error)
	FindByID(id uint) (*models.List, error)
	FindByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardID string, positions []string) error
	GetCardPositions(listPublicID string) ([]uuid.UUID, error)
}

type listRepository struct{}

func NewListRepository() ListRepository {
	return &listRepository{}
}

//	func (r *listRepository) FindByBoardID(boardPublicID string) ([]models.List, error) {
//		var lists []models.List
//		err := config.DB.Joins("JOIN boards ON boards.public_id = lists.board_public_id").
//			Where("boards.public_id = ?", boardPublicID).
//			Order("lists.internal_id ASC").
//			Find(&lists).Error
//		return lists, err
//	}
func (r *listRepository) FindByBoardID(boardPublicID string) ([]models.List, error) {
	var lists []models.List
	err := config.DB.
		Where("board_public_id = ?", boardPublicID).
		Order("internal_id ASC").
		Find(&lists).Error
	return lists, err
}

func (r *listRepository) FindByID(id uint) (*models.List, error) {
	var list models.List
	err := config.DB.First(&list, id).Error
	return &list, err
}

func (r *listRepository) FindByPublicID(publicID string) (*models.List, error) {
	var list models.List
	err := config.DB.Where("public_id = ?", publicID).First(&list).Error
	return &list, err
}

func (r *listRepository) Create(list *models.List) error {
	return config.DB.Create(list).Error
}

func (r *listRepository) Update(list *models.List) error {
	return config.DB.Model(&models.List{}).
		Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title": list.Title,
	}).Error
	return config.DB.Save(list).Error
}

func (r *listRepository) Delete(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

func (r *listRepository) UpdatePositions(boardPublicID string, positions []string) error {
	return config.DB.Model(&models.ListPosition{}).
		Where("board_internal_id = (SELECT internal_id FROM boards WHERE public_id = ?)", boardPublicID).
		Update("list_order", positions).Error
}

func (r *listRepository) GetCardPositions(listPublicID string) ([]uuid.UUID, error) {
	var position models.CardPosition
	err := config.DB.Joins("JOIN lists ON lists.internal_id = card_positions.list_internal_id").
		Where("lists.public_id = ?", listPublicID).
		First(&position).Error
	return position.CardOrder, err
}
