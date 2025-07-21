package usecases

import (
	"encoding/json"
	"fmt"
	"interaction-services/config"
	"interaction-services/domains/timelines"
	"interaction-services/domains/timelines/entities"
	"interaction-services/events"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type timelineUseCase struct {
	repo   timelines.TimelineRepository
	config *config.Config
}

func NewTimelineUseCase(repo timelines.TimelineRepository, config *config.Config) timelines.TimelineUseCase {
	return &timelineUseCase{repo: repo, config: config}
}

func (u *timelineUseCase) getFollowers(userID uuid.UUID) ([]uuid.UUID, error) {
	endpoint := fmt.Sprintf("%s/api/internal/users/%s/followers", u.config.Server.CoreServiceURL, userID)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("service call to '%s' failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("core-service returned non-200 status for followers: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from core-service: %w", err)
	}

	var response struct {
		FollowerIDs []uuid.UUID `json:"follower_ids"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal follower IDs from core-service: %w", err)
	}

	return response.FollowerIDs, nil
}

// AddPostToFollowerTimelines uses the event from the correct package.
func (u *timelineUseCase) AddPostToFollowerTimelines(event *events.PostCreatedEvent) error {
	authorID, _ := uuid.Parse(event.UserID.String())
	postID, _ := uuid.Parse(event.PostID.String())

	followerIDs, err := u.getFollowers(authorID)
	if err != nil {
		log.Printf("CRITICAL: Could not get followers for user %s to fan-out post %s: %v", authorID, postID, err)
		return err
	}

	authorTimelineEntry := &entities.Timeline{
		OwnerID:   authorID,
		PostID:    postID,
		UserID:    authorID,
		ImageURL:  event.ImageURL,
		ThumbURL:  event.ThumbURL,
		Caption:   event.Caption,
		CreatedAt: event.CreatedAt,
	}
	if err := u.repo.AddPostToTimeline(authorTimelineEntry); err != nil {
		log.Printf("ERROR: Failed to add post %s to author's own timeline %s: %v", postID, authorID, err)
	}

	for _, followerID := range followerIDs {
		timelineEntry := &entities.Timeline{
			OwnerID:   followerID,
			PostID:    postID,
			UserID:    authorID,
			ImageURL:  event.ImageURL,
			ThumbURL:  event.ThumbURL,
			Caption:   event.Caption,
			CreatedAt: event.CreatedAt,
		}
		if err := u.repo.AddPostToTimeline(timelineEntry); err != nil {
			log.Printf("ERROR: Failed to add post %s to follower's timeline %s: %v", postID, followerID, err)
		}
	}

	log.Printf("Successfully fanned-out post %s to %d followers.", postID, len(followerIDs))
	return nil
}

func (u *timelineUseCase) GetTimeline(userID uuid.UUID) ([]entities.Timeline, error) {
	return u.repo.GetTimelineForUser(userID)
}
