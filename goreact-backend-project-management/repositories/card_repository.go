package repositories

import (
	"arfian-nurdavicky/project-management/config"
	"arfian-nurdavicky/project-management/models"
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardRepository interface {
	FindByListID(listID string) ([]models.Card, error)
	FindByID(id uint) (*models.Card, error)
	FindCardByPublicID(publicID string) (*models.Card, error)
	Create(card *models.Card) error
	Update(card *models.Card) error
	Delete(id uint) error

	AddLabel(cardID uint, labelID uint) error
	RemoveLabel(cardID uint, labelID uint) error
	AddAssignees(cardID uint, userIDs []uint) error
	RemoveAssignee(cardID uint, userID uint) error

	UpdatePositions(listID string, positions []string) error

	FindByListInternalID(listInternalID int64) (*models.CardPosition, error)
}

type cardRepository struct{}

func NewCardRepository() CardRepository {
	return &cardRepository{}
}

func (r *cardRepository) FindByListID(listPublicID string) ([]models.Card, error) {
	var cards []models.Card
	err := config.DB.Joins("JOIN lists ON lists.internal_id = cards.list_internal_id").
		Where("lists.public_id = ?", listPublicID).
		Order("position ASC").
		Find(&cards).Error
	return cards, err
}

func (r *cardRepository) FindByID(id uint) (*models.Card, error) {
	var card models.Card
	err := config.DB.Preload("Labels").Preload("Assignees").First(&card, id).Error
	return &card, err
}

//func (r *cardRepository) FindCardByPublicID(publicID string) (*models.Card, error) {
//	var card models.Card
//	if err := config.DB.Where("public_id = ?", publicID).First(&card).Error; err != nil {
//		return nil, err
//	}
//	return &card, nil
//}

func (r *cardRepository) Create(card *models.Card) error {
	return config.DB.Create(card).Error
}

func (r *cardRepository) Update(card *models.Card) error {
	return config.DB.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Card{}, id).Error
}

func (r *cardRepository) AddLabel(cardID uint, labelID uint) error {
	cardLabel := models.CardLabel{
		CardID:  int64(cardID),
		LabelID: int64(labelID),
	}
	return config.DB.Create(&cardLabel).Error
}

func (r *cardRepository) RemoveLabel(cardID uint, labelID uint) error {
	return config.DB.Where("card_internal_id = ? AND label_internal_id = ?", cardID, labelID).
		Delete(&models.CardLabel{}).Error
}

func (r *cardRepository) AddAssignees(cardID uint, userIDs []uint) error {
	var assignees []models.CardAssignee
	for _, userID := range userIDs {
		assignees = append(assignees, models.CardAssignee{
			CardID: int64(cardID),
			UserID: int64(userID),
		})
	}

	// gunakan CreateInBatches agar lebih efisien
	//return config.DB.CreateInBatches(&assignees, 100).Error
	// Abaikan jika sudah ada (duplicate)
	return config.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&assignees, 100).Error
}

func (r *cardRepository) RemoveAssignee(cardID uint, userID uint) error {
	return config.DB.Where("card_internal_id = ? AND user_internal_id = ?", cardID, userID).
		Delete(&models.CardAssignee{}).Error
}

func (r *cardRepository) UpdatePositions(listPublicID string, positions []string) error {
	return config.DB.Model(&models.CardPosition{}).
		Where("list_internal_id = (SELECT internal_id FROM lists WHERE public_id = ?)", listPublicID).
		Update("card_order", positions).Error
}

func (r *cardRepository) FindByListInternalID(listInternalID int64) (*models.CardPosition, error) {
	var pos models.CardPosition
	err := config.DB.Where("list_internal_id = ?", listInternalID).First(&pos).Error
	if err != nil {
		return nil, err
	}
	return &pos, nil
}

//func (r *cardRepository) FindCardByPublicID(publicID string) (*models.Card, error) {
//	var card models.Card
//	if err := config.DB.Where("public_id = ?", publicID).First(&card).Error; err != nil {
//		return nil, err
//	}
//	return &card, nil
//}

func (r *cardRepository) FindCardByPublicID(cardPublicID string) (*models.Card, error) {
	var card models.Card
	if err := config.DB.
		Preload("Assignees.User", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("internal_id", "public_id", "name", "email")
		}).
		Preload("Attachments").
		Where("public_id = ?", cardPublicID).
		First(&card).Error; err != nil {
		return nil, err
	}

	// base URL ambil dari ENV (biar fleksibel)
	baseURL := config.AppConfig.AppURL
	if baseURL == "" {
		baseURL = "http://localhost:3030" // fallback
	}

	// generate FileURL untuk setiap attachment
	for i := range card.Attachments {
		card.Attachments[i].FileURL = fmt.Sprintf("%s/files/%s",
			baseURL,
			filepath.Base(card.Attachments[i].File),
		)
	}

	return &card, nil
}
