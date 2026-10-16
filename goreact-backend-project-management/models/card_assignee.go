package models

import "github.com/google/uuid"

//type CardAssignee struct {
//	CardID int64 `json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`
//	UserID int64 `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`
//
//	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:InternalID"`

type CardAssignee struct {
	CardID int64 `json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`
	UserID int64 `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`

	User UserLite `json:"user" gorm:"foreignKey:UserID;references:InternalID"`
}

type UserLite struct {
	InternalID int64     `json:"-" db:"internal_id" gorm:"primaryKey"` // perlu buat relasi, tapi disembunyikan dari JSON
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
}

// kasih tau tabel aslinya tetap "users"
func (UserLite) TableName() string {
	return "users"
}
