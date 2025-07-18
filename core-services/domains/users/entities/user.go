package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Username  string    `gorm:"size:255;not null;"`
	Email     string    `gorm:"size:255;not null;unique"`
	Password  string    `gorm:"size:255;not null;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
