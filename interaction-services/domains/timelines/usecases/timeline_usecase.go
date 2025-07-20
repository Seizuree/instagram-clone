package usecases

import (
	"interaction-services/domains/timelines"
	"interaction-services/domains/timelines/entities"
	"github.com/google/uuid"
	"sort"
)

type timelineUseCase struct {
	repo timelines.TimelineRepository
}

func NewTimelineUseCase(repo timelines.TimelineRepository) timelines.TimelineUseCase {
	return &timelineUseCase{repo: repo}
}

func (u *timelineUseCase) GetTimeline(userID uuid.UUID) ([]entities.TimelineItem, error) {
	followingIDs, err := u.repo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	posts, err := u.repo.GetPostsByUserIDs(followingIDs)
	if err != nil {
		return nil, err
	}

	// Sort by created_at descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}
