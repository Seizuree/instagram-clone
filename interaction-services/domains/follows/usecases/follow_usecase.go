package usecases

import (
	"interaction-services/domains/follows"
	"interaction-services/domains/follows/entities"

	"github.com/google/uuid"
)

type FollowUsecase struct {
	repo follows.FollowRepository
}

func NewFollowUsecase(repo follows.FollowRepository) *FollowUsecase {
	return &FollowUsecase{repo: repo}
}

func (u *FollowUsecase) Follow(followerID, followingID uuid.UUID) error {
	return u.repo.Follow(&entities.Follow{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingID,
	})
}

func (u *FollowUsecase) Unfollow(followerID, followingID uuid.UUID) error {
	return u.repo.Unfollow(followerID, followingID)
}

func (u *FollowUsecase) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	return u.repo.IsFollowing(followerID, followingID)
}

func (u *FollowUsecase) CountFollowers(userID uuid.UUID) (int64, error) {
	return u.repo.CountFollowers(userID)
}

func (u *FollowUsecase) CountFollowing(userID uuid.UUID) (int64, error) {
	return u.repo.CountFollowing(userID)
}
