package postredis

import (
	"time"

	"github.com/google/uuid"
)

type PostModel struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	PostedAt  time.Time `json:"posted_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
