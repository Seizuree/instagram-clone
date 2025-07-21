package repositories

import (
	"core-services/domains/follow"
	"core-services/domains/follow/entities"
	"core-services/infrastructures"
	"errors"

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
	return f.db.GetInstance().
		Where(entities.Follow{FollowerID: followerID, FollowingID: followingID}).
		FirstOrCreate(&entities.Follow{
			ID:          uuid.New(),
			FollowerID:  followerID,
			FollowingID: followingID,
		}).Error
}

// Unfollow implements follow.FollowRepository.
func (f *followRepository) Unfollow(followerID uuid.UUID, followingID uuid.UUID) error {
	// Perform the delete operation
	result := f.db.GetInstance().Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&entities.Follow{})

	// Check for any database errors first
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were actually deleted
	if result.RowsAffected == 0 {
		return errors.New("user already unfollowed or was never followed")
	}

	return nil
}

// GetFollowerIDs implements follow.FollowRepository.
func (f *followRepository) GetFollowerIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error) {
	var followerIDs []uuid.UUID
	err := f.db.GetInstance().
		Model(&entities.Follow{}).
		Where("following_id = ?", userID).       // Find all records where the user is the one being followed.
		Pluck("follower_id", &followerIDs).Error // Extract only the follower_id column into the slice.

	return followerIDs, err
}

func (f *followRepository) GetFollowCounts(userID uuid.UUID) (followerCount int64, followingCount int64, err error) {
	err = f.db.GetInstance().Model(&entities.Follow{}).Where("following_id = ?", userID).Count(&followerCount).Error
	if err != nil {
		return 0, 0, err
	}

	err = f.db.GetInstance().Model(&entities.Follow{}).Where("follower_id = ?", userID).Count(&followingCount).Error
	if err != nil {
		return 0, 0, err
	}

	return followerCount, followingCount, nil
}
