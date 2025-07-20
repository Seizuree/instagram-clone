package repositories

import (
	"interaction-services/domains/follows"
	"interaction-services/domains/follows/entities"
	"interaction-services/infrastructures"

	"github.com/google/uuid"
)

type followRepository struct {
	db infrastructures.Database
}

func NewFollowRepository(db infrastructures.Database) follows.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(f *entities.Follow) error {
	return r.db.GetInstance().Create(f).Error
}

func (r *followRepository) Unfollow(followerID, followingID uuid.UUID) error {
	return r.db.GetInstance().Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entities.Follow{}).Error
}

func (r *followRepository) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.GetInstance().Model(&entities.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

func (r *followRepository) CountFollowers(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.GetInstance().Model(&entities.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *followRepository) CountFollowing(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.GetInstance().Model(&entities.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}
