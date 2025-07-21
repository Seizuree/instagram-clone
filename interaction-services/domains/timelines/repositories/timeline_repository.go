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

// GetTimelineForUser implements timelines.TimelineRepository.
func (t *timelineRepo) GetTimelineForUser(ownerID uuid.UUID) ([]entities.Timeline, error) {
	var timeline []entities.Timeline

	err := t.db.GetInstance().
		Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Limit(50).
		Find(&timeline).Error
	return timeline, err
}

// AddPostToTimeline implements timelines.TimelineRepository.
func (t *timelineRepo) AddPostToTimeline(timelineEntry *entities.Timeline) error {
	timelineEntry.ID = uuid.New()
	return t.db.GetInstance().Create(timelineEntry).Error
}

func (t *timelineRepo) AddPostsToTimeline(timelineEntries []*entities.Timeline) error {
	for i := range timelineEntries {
		timelineEntries[i].ID = uuid.New()
	}
	return t.db.GetInstance().CreateInBatches(timelineEntries, 100).Error
}
