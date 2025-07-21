package repositories

import (
	"errors"
	"interaction-services/domains/likes"
	"interaction-services/domains/likes/entities"
	"interaction-services/infrastructures"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type likeRepository struct {
	db infrastructures.Database
}

func NewLikeRepository(db infrastructures.Database) likes.LikeRepository {
	return &likeRepository{db: db}
}

func (l *likeRepository) LikePost(userID, postID uuid.UUID) error {
	var existingLike entities.Like
	err := l.db.GetInstance().Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record is not found, it means the user hasn't liked the post yet.
			// So, we can proceed to create a new like.
			newLike := &entities.Like{
				ID:     uuid.New(),
				UserID: userID,
				PostID: postID,
			}
			return l.db.GetInstance().Create(newLike).Error
		}
		// If there is another type of database error, we should return it.
		return err
	}

	// If err is nil, it means the user has already liked the post.
	return errors.New("post already liked")
}

func (l *likeRepository) UnlikePost(userID, postID uuid.UUID) error {
	result := l.db.GetInstance().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&entities.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("post was not liked")
	}

	return nil
}

func (l *likeRepository) CountLikes(postID uuid.UUID) (int64, error) {
	var count int64
	err := l.db.GetInstance().Model(&entities.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
