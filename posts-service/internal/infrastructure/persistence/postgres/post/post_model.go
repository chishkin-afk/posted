package postpg

import (
	"time"

	"github.com/google/uuid"
)

type PostModel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	OwnerID   uuid.UUID `gorm:"not null;index"`
	Title     string    `gorm:"not null"`
	Body      string    `gorm:"not null"`
	PostedAt  time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
