package follow

import "github.com/google/uuid"

type FollowUseCase interface {
	Follow(followerID uuid.UUID, username string) error
	Unfollow(followerID uuid.UUID, username string) error
}

type FollowRepository interface {
	Follow(followerID, followingID uuid.UUID) error
	Unfollow(followerID, followingID uuid.UUID) error
}
