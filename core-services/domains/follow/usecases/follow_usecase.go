package usecases

import (
	"core-services/domains/follow"
	"core-services/domains/users"

	"github.com/google/uuid"
)

type followUseCase struct {
	followRepo follow.FollowRepository
	userRepo   users.UserRepository
}

func NewFollowUseCase(followRepo follow.FollowRepository, userRepo users.UserRepository) follow.FollowUseCase {
	return &followUseCase{followRepo: followRepo, userRepo: userRepo}
}

// Follow implements follow.FollowUsecase.
func (f *followUseCase) Follow(followerID uuid.UUID, username string) error {
	userToFollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Follow(followerID, userToFollow.ID)
}

// Unfollow implements follow.FollowUsecase.
func (f *followUseCase) Unfollow(followerID uuid.UUID, username string) error {
	userToUnfollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Unfollow(followerID, userToUnfollow.ID)
}
