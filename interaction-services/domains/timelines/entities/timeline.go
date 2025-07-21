package entities

import (
	"time"

	"github.com/google/uuid"
)

type Timeline struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID   uuid.UUID `gorm:"not null;index"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	ImageURL  string    `json:"image_url"`
	ThumbURL  string    `json:"thumb_url"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"created_at"`
}
