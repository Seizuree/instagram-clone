package timelines

import (
	"interaction-services/domains/timelines/entities"
	"interaction-services/events"

	"github.com/google/uuid"
)

type TimelineUseCase interface {
	GetTimeline(userID uuid.UUID) ([]entities.Timeline, error)
	AddPostToFollowerTimelines(event *events.PostCreatedEvent) error
	AddPostsToFollowerTimeline(userID uuid.UUID, posts []*entities.Timeline) error
}

type TimelineRepository interface {
	GetTimelineForUser(ownerID uuid.UUID) ([]entities.Timeline, error)
	AddPostToTimeline(timelineEntry *entities.Timeline) error
	AddPostsToTimeline(timelineEntries []*entities.Timeline) error
}
