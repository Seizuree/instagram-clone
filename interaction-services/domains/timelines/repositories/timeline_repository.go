package repositories

import (
	"interaction-services/domains/timelines"
	"interaction-services/domains/timelines/entities"
	"interaction-services/infrastructures"
	"github.com/google/uuid"
)

type timelineRepo struct {
	db infrastructures.Database
}

func NewTimelineRepository(db infrastructures.Database) timelines.TimelineRepository {
	return &timelineRepo{db: db}
}

func (r *timelineRepo) GetFollowing(userID uuid.UUID) ([]uuid.UUID, error) {
	var followIDs []uuid.UUID
	err := r.db.GetInstance().Raw(`
		SELECT following_id FROM follows WHERE follower_id = ?
	`, userID).Scan(&followIDs).Error
	return followIDs, err
}

func (r *timelineRepo) GetPostsByUserIDs(userIDs []uuid.UUID) ([]entities.TimelineItem, error) {
	if len(userIDs) == 0 {
		return []entities.TimelineItem{}, nil
	}

	var posts []entities.TimelineItem
	err := r.db.GetInstance().Raw(`
		SELECT id as post_id, user_id, image_url, thumb_url, caption, created_at
		FROM posts WHERE user_id IN (?) ORDER BY created_at DESC
	`, userIDs).Scan(&posts).Error
	return posts, err
}
