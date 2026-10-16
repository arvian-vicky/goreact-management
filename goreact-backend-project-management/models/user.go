package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

//	type User struct {
//		InternalID uint           `json:"internal_id" gorm:"primaryKey"`
//		PublicID   uuid.UUID      `json:"public_id" db:"public_id"`
//		Name       string         `json:"name"`
//		Email      string         `json:"email" gorm:"unique"`
//		Password   string         `json:"password,omitempty" gorm:"column:password" `
//		Role       string         `json:"role"` // e.g. "admin", "user"
//		CreatedAt  time.Time      `json:"created_at"`
//		UpdatedAt  time.Time      `json:"updated_at"`
//		DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
//	}
type User struct {
	InternalID int64          `json:"internal_id" db:"internal_id" gorm:"primaryKey"`
	PublicID   uuid.UUID      `json:"public_id" db:"public_id"`
	Name       string         `json:"name" db:"name"`
	Email      string         `json:"email" db:"email" gorm:"unique"`
	Password   string         `json:"password,omitempty"  db:"password" gorm:"column:password" `
	Role       string         `json:"role"` // e.g. "admin", "user"
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserResponse struct {
	PublicID  uuid.UUID `json:"public_id" db:"public_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
