package services

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"arfian-nurdavicky/project-management/models/types"
	"arfian-nurdavicky/project-management/repositories"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)

	AddLabel(cardPublicID, labelPublicID string) error
	RemoveLabel(cardPublicID, labelPublicID string) error

	AddAssignees(cardPublicID string, userPublicIDs []string) error
	RemoveAssignee(cardPublicID, userPublicID string) error

	UpdatePositions(listPublicID string, positions []string) error
	UpdateCardList(card *models.Card, listPublicID string) error

	GetCardDetail(cardPublicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo  repositories.CardRepository
	listRepo  repositories.ListRepository
	labelRepo repositories.LabelRepository
	userRepo  repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	labelRepo repositories.LabelRepository,
	userRepo repositories.UserRepository,
) CardService {
	return &cardService{cardRepo, listRepo, labelRepo, userRepo}
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	// 1. Verifikasi list ada
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("list not found: %w", err)
	}

	// 2. Ambil posisi card
	position, err := s.cardRepo.FindByListInternalID(list.InternalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get card positions: %w", err)
	}

	// 3. Ambil semua cards dari list
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	// 4. Jika ada posisi → urutkan cards sesuai card_order
	if position != nil && len(position.CardOrder) > 0 {
		cards = sortCardsByPosition(cards, position.CardOrder)
	}

	return cards, nil
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindCardByPublicID(publicID)
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// 1. Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Set list_internal_id
	card.ListID = list.InternalID

	// 3. Generate public_id jika belum ada
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	// 4. Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 5. Simpan card
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	// 6. Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru jika belum ada
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}

			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get card position: %w", err)
		}
	} else {
		// Tambahkan card baru ke urutan
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position: %w", err)
		}
	}

	// 7. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	// 1. Ambil card lama
	existingCard, err := s.cardRepo.FindCardByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	// 2. Ambil list baru
	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 3. Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 4. Jika pindah list → hapus dari posisi list lama & tambah ke list baru
	if existingCard.ListID != newList.InternalID {
		// Hapus dari list lama
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err == nil {
			// buat filtered slice tanpa existingCard.PublicID
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}
			// Update menggunakan models.UUIDArray(filtered) agar Value() dipanggil
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update old card position: %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("failed to get old card position: %w", err)
		}

		// Tambah ke list baru
		var newPos models.CardPosition
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}
			if err := tx.Create(&newPos).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list: %w", err)
			}
		} else if res.Error == nil {
			// append dengan models.UUIDArray agar Value() dipanggil
			updatedOrder := append(newPos.CardOrder, existingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(updatedOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update new card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get new card position: %w", res.Error)
		}
	}

	// 5. Update data card (pastikan ID tidak berubah)
	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = newList.InternalID

	// Gunakan Save (atau Updates jika ingin only non-zero fields) — disini kita Save agar mirip sebelumnya
	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card: %w", err)
	}

	// 6. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) AddLabel(cardPublicID, labelPublicID string) error {
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	label, err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("label not found")
	}

	return s.cardRepo.AddLabel(uint(card.InternalID), uint(label.InternalID))
}

func (s *cardService) RemoveLabel(cardPublicID, labelPublicID string) error {
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	label, err := s.labelRepo.FindByPublicID(labelPublicID)
	if err != nil {
		return errors.New("label not found")
	}

	return s.cardRepo.RemoveLabel(uint(card.InternalID), uint(label.InternalID))
}

func (s *cardService) AddAssignees(cardPublicID string, userPublicIDs []string) error {
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	var userIDs []uint
	for _, userPublicID := range userPublicIDs {
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil {
			return fmt.Errorf("user %s not found", userPublicID)
		}
		userIDs = append(userIDs, uint(user.InternalID))
	}

	return s.cardRepo.AddAssignees(uint(card.InternalID), userIDs)
}

func (s *cardService) RemoveAssignee(cardPublicID, userPublicID string) error {
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return errors.New("card not found")
	}

	user, err := s.userRepo.FindByPublicID(userPublicID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.cardRepo.RemoveAssignee(uint(card.InternalID), uint(user.InternalID))
}

func (s *cardService) UpdatePositions(listPublicID string, positions []string) error {
	// 1. Pastikan list ada
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Convert []string → UUIDArray
	var cardOrder types.UUIDArray
	for _, p := range positions {
		id, err := uuid.Parse(p)
		if err != nil {
			return fmt.Errorf("invalid card id: %w", err)
		}
		cardOrder = append(cardOrder, id)
	}

	// 3. Mulai transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 4. Update card_positions
	var position models.CardPosition
	res := tx.Where("list_internal_id = ?", list.InternalID).First(&position)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Kalau belum ada, bikin baru
		position = models.CardPosition{
			PublicID:  uuid.New(),
			ListID:    list.InternalID,
			CardOrder: cardOrder,
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create card position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get card position: %w", res.Error)
	} else {
		// Kalau sudah ada, update array-nya
		if err := tx.Model(&position).
			Where("internal_id = ?", position.InternalID).
			Update("card_order", cardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card positions: %w", err)
		}
	}

	// 5. Commit
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) UpdateCardList(card *models.Card, listPublicID string) error {
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}
	card.ListID = list.InternalID
	return nil
}

func sortCardsByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	// 1. Buat map untuk pencarian cepat O(1)
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i // Simpan posisi (0, 1, 2, ...)
	}

	// 2. Tentukan "indeks default" untuk item yang TIDAK ada di 'order'
	// Kita gunakan angka yang lebih besar dari indeks manapun di 'order'.
	// len(order) adalah pilihan yang sempurna.
	defaultIndex := len(order)

	// 3. Sortir slice
	sort.SliceStable(cards, func(i, j int) bool {
		// Dapatkan indeks untuk card[i], jika tidak ada, gunakan defaultIndex
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}

		// Dapatkan indeks untuk card[j], jika tidak ada, gunakan defaultIndex
		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}

		// Jika kedua indeks sama (keduanya ada di order dengan posisi sama,
		// atau KEDUANYA TIDAK ADA di order), maka urutkan berdasarkan CreatedAt.
		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}

		// Jika tidak, urutkan berdasarkan indeks.
		// Item yang ada di 'order' (cth: idxI = 1) akan selalu
		// lebih kecil dari item yang tidak ada (cth: idxJ = defaultIndex).
		return idxI < idxJ
	})

	return cards
}

//	func (s *cardService) GetCardDetail(cardPublicID string) (*models.Card, error) {
//		card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
//		if err != nil {
//			return nil, err
//		}
//		return card, nil
//	}
func (s *cardService) GetCardDetail(cardPublicID string) (*models.Card, error) {
	// Ambil card dari repository, preload assignees dan attachments
	card, err := s.cardRepo.FindCardByPublicID(cardPublicID)
	if err != nil {
		return nil, err
	}
	return card, nil
}
