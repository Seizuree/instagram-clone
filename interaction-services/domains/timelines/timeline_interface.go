package timelines

import (
	"interaction-services/domains/timelines/entities"
	"github.com/google/uuid"
)

type TimelineUseCase interface {
	GetTimeline(userID uuid.UUID) ([]entities.TimelineItem, error)
}

type TimelineRepository interface {
	GetFollowing(userID uuid.UUID) ([]uuid.UUID, error)
	GetPostsByUserIDs(userIDs []uuid.UUID) ([]entities.TimelineItem, error)
}
