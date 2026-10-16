package models

import (
	"arfian-nurdavicky/project-management/models/types"

	"github.com/google/uuid"
)

type ListPosition struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id" gorm:"public_id"`
	BoardID    int64     `json:"board_internal_id" db:"board_internal_id" gorm:"column:board_internal_id"`
	//ListOrder  []uuid.UUID `json:"list_order" db:"list_order" pg:"type:uuid[]"` // Tetap sebagai []uuid.UUID
	ListOrder types.UUIDArray `json:"list_order" gorm:"type:uuid[]"`
}

func (ListPosition) TableName() string {
	return "list_positions"
}
