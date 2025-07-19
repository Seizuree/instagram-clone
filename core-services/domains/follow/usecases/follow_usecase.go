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
func (f *followUseCase) Follow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)
	if err != nil {
		return err
	}

	userToFollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Follow(parsedFollowerID, userToFollow.ID)
}

// Unfollow implements follow.FollowUsecase.
func (f *followUseCase) Unfollow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)

	if err != nil {
		return err
	}

	userToUnfollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Unfollow(parsedFollowerID, userToUnfollow.ID)
}
