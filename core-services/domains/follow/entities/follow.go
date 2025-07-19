package entities

import (
	"core-services/domains/users/entities"
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID     `gorm:"primaryKey"`
	FollowerID  uuid.UUID     `gorm:"not null"`
	Follower    entities.User `gorm:"foreignKey:FollowerID;references:ID"`
	FollowingID uuid.UUID     `gorm:"not null"`
	Following   entities.User `gorm:"foreignKey:FollowingID;references:ID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
