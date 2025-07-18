package repositories

import (
	"core-services/domains/follow"
	"core-services/domains/follow/entities"
	"core-services/infrastructures"

	"github.com/google/uuid"
)

type followRepository struct {
	db infrastructures.Database
}

func NewFollowRepository(db infrastructures.Database) follow.FollowRepository {
	return &followRepository{db: db}
}

// Follow implements follow.FollowRepository.
func (f *followRepository) Follow(followerID uuid.UUID, followingID uuid.UUID) error {
	follow := &entities.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return f.db.GetInstance().Create(follow).Error
}

// Unfollow implements follow.FollowRepository.
func (f *followRepository) Unfollow(followerID uuid.UUID, followingID uuid.UUID) error {
	return f.db.GetInstance().Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entities.Follow{}).Error
}
