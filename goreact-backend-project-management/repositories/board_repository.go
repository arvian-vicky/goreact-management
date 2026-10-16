package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"time"

	"github.com/google/uuid"
)

type BoardRepository interface {
	FindAll() ([]models.Board, error)
	FindAllByUser(userID string) ([]models.Board, error)
	FindByID(id uint) (*models.Board, error)
	FindByPublicID(publicID string) (*models.Board, error)
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(id uint) error
	AddMember(boardID uint, userID uint) error
	AddMembers(boardID uint, userIDs []uint) error
	RemoveMember(boardID uint, userID uint) error
	RemoveMembers(boardID uint, userIDs []uint) error
	GetListPositions(boardPublicID string) ([]uuid.UUID, error)
	FindAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error)
	FindAllByUserPaginate(userID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
}

type boardRepository struct{}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) FindAll() ([]models.Board, error) {
	var boards []models.Board
	err := config.DB.Find(&boards).Error
	return boards, err
}

func (r *boardRepository) FindAllByUser(userPublicID string) ([]models.Board, error) {
	var boards []models.Board
	err := config.DB.
		Where("owner_public_id = ? OR internal_id IN (SELECT board_members.board_internal_id FROM board_members JOIN users ON users.internal_id = board_members.user_internal_id WHERE users.public_id = ?)", userPublicID, userPublicID).
		Find(&boards).Error
	return boards, err
}

func (r *boardRepository) FindByID(id uint) (*models.Board, error) {
	var board models.Board
	err := config.DB.First(&board, id).Error
	return &board, err
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board
	err := config.DB.Where("public_id = ?", publicID).First(&board).Error
	return &board, err
}

func (r *boardRepository) Create(board *models.Board) error {

	return config.DB.Create(board).Error
}

func (r *boardRepository) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).
		Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error

	//return config.DB.Save(board).Error
}

func (r *boardRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Board{}, id).Error
}

func (r *boardRepository) AddMember(boardID uint, userID uint) error {
	member := models.BoardMember{
		BoardID:  int64(boardID),
		UserID:   int64(userID),
		JoinedAt: time.Now(),
	}
	return config.DB.Create(&member).Error
}

func (r *boardRepository) AddMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	now := time.Now()
	var members []models.BoardMember

	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardID:  int64(boardID),
			UserID:   int64(userID),
			JoinedAt: now,
		})
	}

	return config.DB.Create(&members).Error
}

func (r *boardRepository) RemoveMember(boardID uint, userID uint) error {
	return config.DB.Where("board_internal_id = ? AND user_internal_id = ?", boardID, userID).
		Delete(&models.BoardMember{}).Error
}

func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return nil
	}

	return config.DB.
		Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).
		Delete(&models.BoardMember{}).Error
}

// Ubah implementasi
func (r *boardRepository) GetListPositions(boardPublicID string) ([]uuid.UUID, error) {
	var position models.ListPosition
	err := config.DB.Joins("JOIN boards ON boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).
		First(&position).Error
	return position.ListOrder, err
}

func (r *boardRepository) FindAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var boards []models.Board
	var total int64

	query := config.DB.Model(&models.Board{})

	// Filter (misalnya berdasarkan title)
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}

	// Hitung total dulu sebelum limit+offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at desc")
	}

	// Pagination
	if err := query.Limit(limit).Offset(offset).Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) FindAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var boards []models.Board
	var total int64

	// Base query dengan filter user ownership atau membership
	query := config.DB.Model(&models.Board{}).
		Where("owner_public_id = ? OR internal_id IN ("+
			"SELECT board_members.board_internal_id FROM board_members "+
			"JOIN users ON users.internal_id = board_members.user_internal_id "+
			"WHERE users.public_id = ?)", userPublicID, userPublicID)

	// Filter tambahan (misalnya berdasarkan title)
	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}

	// Hitung total dulu sebelum limit dan offset
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at desc")
	}

	// Pagination
	if err := query.Limit(limit).Offset(offset).Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}
