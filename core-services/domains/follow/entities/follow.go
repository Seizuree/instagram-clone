package entities

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
	CreatedAt   time.Time
}
