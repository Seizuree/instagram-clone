package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    ImageURL  string    `gorm:"not null" json:"image_url"`
    ThumbURL  string    `gorm:"not null" json:"thumb_url"`
    Caption   string    `gorm:"type:text" json:"caption"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
