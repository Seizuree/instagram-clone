package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID    uuid.UUID `gorm:"not null"`
	ImageURL  string    `gorm:"size:255;not null;"`
	Caption   string    `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
