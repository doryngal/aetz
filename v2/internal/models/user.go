package models

import (
	"time"
)

type User struct {
	ID           int64  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string `json:"name" gorm:"type:varchar(50);not null" validate:"required"`
	Email        string `json:"email" gorm:"type:varchar(100);unique;not null" validate:"required,email"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null" validate:"required"`
	Role         string `json:"role" gorm:"type:varchar(20);default:'user'" validate:"required"`
	IsActive     bool   `json:"is_active" gorm:"type:boolean;not null" validate:"required"`

	CompanyID int64 `json:"company_id" gorm:"type:uuid;not null" validate:"required"`

	ResetCode          *string    `json:"-" gorm:"type:varchar(4);default:null"` // Код восстановления
	ResetCodeExpiresAt *time.Time `json:"-" gorm:"default:null"`                 // Время истечения кода
	LastResetSentAt    *time.Time `json:"-" gorm:"default:null"`                 // Время последней отправки кода

	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}
