package follows

import (
	"interaction-services/domains/follows/entities"
	"github.com/google/uuid"
)

type FollowRepository interface {
	Follow(follow *entities.Follow) error
	Unfollow(followerID, followingID uuid.UUID) error
	IsFollowing(followerID, followingID uuid.UUID) (bool, error)
	CountFollowers(userID uuid.UUID) (int64, error)
	CountFollowing(userID uuid.UUID) (int64, error)
}
