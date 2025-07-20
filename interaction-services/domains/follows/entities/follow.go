package entities

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FollowerID uuid.UUID `gorm:"type:uuid;not null"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
}
