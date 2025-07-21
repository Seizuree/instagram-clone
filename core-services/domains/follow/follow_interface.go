package follow

import "github.com/google/uuid"

type FollowUseCase interface {
	Follow(followerID string, username string) error
	Unfollow(followerID string, username string) error
	GetFollowerIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error)
}

type FollowRepository interface {
	Follow(followerID, followingID uuid.UUID) error
	Unfollow(followerID, followingID uuid.UUID) error
	GetFollowerIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error)
	GetFollowCounts(userID uuid.UUID) (followerCount int64, followingCount int64, err error)
}
