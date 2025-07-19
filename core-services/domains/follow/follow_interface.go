package follow

import "github.com/google/uuid"

type FollowUseCase interface {
	Follow(followerID string, username string) error
	Unfollow(followerID string, username string) error
}

type FollowRepository interface {
	Follow(followerID, followingID uuid.UUID) error
	Unfollow(followerID, followingID uuid.UUID) error
}
