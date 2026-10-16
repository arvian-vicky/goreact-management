package services

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/models/types"
	"arfian-nurdavicky/project-management/repositories"
	"arfian-nurdavicky/project-management/utils"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListWithOrder struct {
	Positions []uuid.UUID
	Lists     []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardPublicID string, positions []uuid.UUID) error
}

type listService struct {
	listRepo    repositories.ListRepository
	boardRepo   repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository
}

func NewListService(
	listRepo repositories.ListRepository,
	boardRepo repositories.BoardRepository,
	listPosRepo repositories.ListPositionRepository,
) ListService {
	return &listService{listRepo, boardRepo, listPosRepo}
}

func (s *listService) GetByBoardID(boardPublicID string) (*ListWithOrder, error) {
	// Verifikasi board ada
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}
	position, err := s.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list order: " + err.Error())
	}

	lists, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get lists: " + err.Error())
	}

	orderedLists := utils.SortListsByPosition(lists, position)

	return &ListWithOrder{
		Positions: position,
		Lists:     orderedLists,
	}, nil
}

func (s *listService) GetByID(id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicID)
}

//	func (s *listService) Create(list *models.List) error {
//		// Pastikan board ada
//		fmt.Println(list.BoardPublicID.String())
//		_, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
//		if err != nil {
//			return errors.New("board not found")
//		}
//		list.PublicID = uuid.New()
//		return s.listRepo.Create(list)
//	}
func (s *listService) Create(list *models.List) error {
	// 1. Validasi board
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return fmt.Errorf("failed to get board: %w", err)
	}

	// Set internal board ID
	list.BoardInternalID = board.InternalID

	// Generate public ID jika belum ada
	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	// Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 2. Simpan list baru
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list: %w", err)
	}

	// 3. Update list_position
	var position models.ListPosition

	// Coba ambil list_positions untuk board ini
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Buat baru jika belum ada
		position = models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get list position: %w", res.Error)
	} else {
		// Tambahkan ID baru
		position.ListOrder = append(position.ListOrder, list.PublicID)
		// Update ke DB
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update list position: %w", err)
		}
	}

	// 4. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (s *listService) Update(list *models.List) error {
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePositions(boardPublicID string, positions []uuid.UUID) error {
	// Verifikasi board ada
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// Dapatkan list position
	position, err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("list position not found")
	}

	// Update list order
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}
