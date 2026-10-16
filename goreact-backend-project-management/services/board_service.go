package services

import (
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/repositories"
	"errors"

	"github.com/google/uuid"
)

type BoardService interface {
	GetAll() ([]models.Board, error)
	GetAllByUser(userPublicID string) ([]models.Board, error)
	GetByID(id uint) (*models.Board, error)
	GetByPublicID(publicID string) (*models.Board, error)
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(id uint) error
	AddMember(boardPublicID, userPublicID string) error
	AddMembers(boardPublicID string, userPublicIDs []string) error
	RemoveMember(boardPublicID, userPublicID string) error
	RemoveMembers(boardPublicID string, userPublicIDs []string) error
	IsMember(boardPublicID, userPublicID string) (bool, error)
	GetMembers(publicID string) ([]models.User, error)
	GetAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error)
	GetAllByUserPaginate(userID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
}

type boardService struct {
	boardRepo       repositories.BoardRepository
	userRepo        repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
}

func (s *boardService) GetAll() ([]models.Board, error) {
	return s.boardRepo.FindAll()
}

func (s *boardService) GetAllByUser(userPublicID string) ([]models.Board, error) {
	return s.boardRepo.FindAllByUser(userPublicID)
}

func (s *boardService) GetByID(id uint) (*models.Board, error) {
	return s.boardRepo.FindByID(id)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) Create(board *models.Board) error {
	// Pastikan user yang membuat board ada
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())

	if err != nil {
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) Delete(id uint) error {
	return s.boardRepo.Delete(id)
}

func (s *boardService) AddMember(boardPublicID, userPublicID string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.boardRepo.AddMember(uint(board.InternalID), uint(user.InternalID))
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// Validasi semua user dan kumpulkan internal IDs
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return errors.New("user not found: " + userPublicID)
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	// Cek keanggotaan sekaligus untuk semua user
	existingMembers, err := s.boardMemberRepo.GetMembers(board.PublicID.String())
	if err != nil {
		return err
	}

	// Buat map untuk pengecekan cepat
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true
	}

	// Siapkan hanya user yang belum menjadi member
	var newMemberIDs []uint
	for _, userID := range userInternalIDs {
		if !memberMap[userID] {
			newMemberIDs = append(newMemberIDs, userID)
		}
	}

	if len(newMemberIDs) == 0 {
		return nil // Semua user sudah menjadi member
	}

	return s.boardRepo.AddMembers(uint(board.InternalID), newMemberIDs)
}

func (s *boardService) RemoveMember(boardPublicID, userPublicID string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.boardRepo.RemoveMember(uint(board.InternalID), uint(user.InternalID))
}

// Implementasi method
func (s *boardService) RemoveMembers(boardPublicID string, userPublicIDs []string) error {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	// Validasi semua user
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			//return errors.New("user not found: " + userPublicID)
			continue
		}
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	//// Dapatkan semua member yang ada
	//existingMembers, err := s.boardRepo.GetMembers(uint(board.InternalID))
	//if err != nil {
	//	return err
	//}
	// Cek keanggotaan sekaligus untuk semua user
	existingMembers, err := s.boardMemberRepo.GetMembers(board.PublicID.String())
	if err != nil {
		return err
	}

	// Buat map untuk pengecekan cepat
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers {
		memberMap[uint(member.InternalID)] = true
	}

	// Siapkan hanya user yang memang member
	var membersToRemove []uint
	for _, userID := range userInternalIDs {
		if memberMap[userID] {
			membersToRemove = append(membersToRemove, userID)
		}
	}

	if len(membersToRemove) == 0 {
		return nil // Tidak ada member yang perlu dihapus
	}

	return s.boardRepo.RemoveMembers(uint(board.InternalID), membersToRemove)
}

func (s *boardService) IsMember(boardPublicID, userPublicID string) (bool, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return false, errors.New("board not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return false, errors.New("user not found")
	}

	//return s.boardRepo.IsMember(uint(board.InternalID), uint(user.InternalID))
	return s.boardMemberRepo.IsMember(string(board.InternalID), string(user.InternalID))
}

func (s *boardService) GetMembers(boardPublicID string) ([]models.User, error) {
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, err
	}
	return s.boardMemberRepo.GetMembers(board.PublicID.String())
}

func (s *boardService) GetAllPaginate(filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllPaginate(filter, sort, limit, offset)
}

func (s *boardService) GetAllByUserPaginate(userID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	return s.boardRepo.FindAllByUserPaginate(userID, filter, sort, limit, offset)
}
