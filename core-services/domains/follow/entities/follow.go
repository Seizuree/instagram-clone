package entities

import (
	"core-services/domains/users/entities"
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID     `gorm:"primaryKey"`
	FollowerID  uuid.UUID     `gorm:"not null;uniqueIndex:idx_follower_following"`
	Follower    entities.User `gorm:"foreignKey:FollowerID;references:ID;onDelete:CASCADE"`
	FollowingID uuid.UUID     `gorm:"not null;uniqueIndex:idx_following_follower"`
	Following   entities.User `gorm:"foreignKey:FollowingID;references:ID;onDelete:CASCADE"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
