package repositories

import (
	"interaction-services/domains/likes"
	"interaction-services/domains/likes/entities"
	"interaction-services/infrastructures"

	"github.com/google/uuid"
)

type likeRepository struct {
	db infrastructures.Database
}

func NewLikeRepository(db infrastructures.Database) likes.LikeRepository {
	return &likeRepository{db: db}
}

func (l *likeRepository) LikePost(userID, postID uuid.UUID) error {
	like := &entities.Like{
		ID:     uuid.New(),
		UserID: userID,
		PostID: postID,
	}
	return l.db.GetInstance().Create(like).Error
}

func (l *likeRepository) UnlikePost(userID, postID uuid.UUID) error {
	return l.db.GetInstance().Where("user_id = ? AND post_id = ?", userID, postID).Delete(&entities.Like{}).Error
}
