package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	PostID    uuid.UUID `gorm:"not null"`
	UserID    uuid.UUID `gorm:"not null"`
	Comment   string    `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
