package userpg

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Email        string    `gorm:"not null;indexUnique"`
	PasswordHash string    `gorm:"not null"`
	Nickname     string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}
