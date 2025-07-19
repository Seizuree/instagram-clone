package entities

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	PostID    uuid.UUID `gorm:"not null"`
	UserID    uuid.UUID `gorm:"not null"`
	CreatedAt time.Time
}
