package usecases

import (
	"bytes"
	"context"
	"core-services/config"
	"core-services/domains/follow"
	"core-services/domains/users"
	"core-services/events"
	"core-services/infrastructures"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type followUseCase struct {
	followRepo follow.FollowRepository
	userRepo   users.UserRepository
	rabbitMQ   *infrastructures.RabbitMQ
	config     *config.Config
}

func NewFollowUseCase(followRepo follow.FollowRepository, userRepo users.UserRepository, rabbitMQ *infrastructures.RabbitMQ, config *config.Config) follow.FollowUseCase {
	return &followUseCase{followRepo: followRepo, userRepo: userRepo, rabbitMQ: rabbitMQ, config: config}
}

func (f *followUseCase) populateTimeline(followerID, followedID uuid.UUID) error {
	// Define structs for data transformation
	type PostFromService struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		ImageURL  string    `json:"image_url"`
		ThumbURL  string    `json:"thumb_url"`
		Caption   string    `json:"caption"`
		CreatedAt time.Time `json:"created_at"`
	}

	type TimelinePostPayload struct {
		PostID    uuid.UUID `json:"post_id"`
		UserID    uuid.UUID `json:"user_id"`
		ImageURL  string    `json:"image_url"`
		ThumbURL  string    `json:"thumb_url"`
		Caption   string    `json:"caption"`
		CreatedAt time.Time `json:"created_at"`
	}

	// 1. Get posts of the followed user from post-services
	postServiceURL := f.config.Server.PostServiceURL
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/posts/user/%s", postServiceURL, followedID.String()), nil)
	if err != nil {
		return fmt.Errorf("failed to create request to post-service: %w", err)
	}
	req.Header.Set("X-User-ID", followerID.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get posts from post-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post-service returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from post-service: %w", err)
	}

	// 2. Decode the response from post-service
	var postsFromService []PostFromService
	if err := json.Unmarshal(body, &postsFromService); err != nil {
		return fmt.Errorf("failed to unmarshal posts from post-service: %w", err)
	}

	// 3. Transform the data into the format expected by interaction-service
	var timelinePayload []TimelinePostPayload
	for _, post := range postsFromService {
		timelinePayload = append(timelinePayload, TimelinePostPayload{
			PostID:    post.ID, // <-- The critical mapping step
			UserID:    post.UserID,
			ImageURL:  post.ImageURL,
			ThumbURL:  post.ThumbURL,
			Caption:   post.Caption,
			CreatedAt: post.CreatedAt,
		})
	}

	// 4. Marshal the new, correct payload
	payloadBytes, err := json.Marshal(timelinePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal timeline payload: %w", err)
	}

	// 5. Send the transformed posts to interaction-services to update the timeline
	interactionServiceURL := f.config.Server.InteractionServiceURL
	ireq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/internal/timelines", interactionServiceURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request to interaction-service: %w", err)
	}
	ireq.Header.Set("Content-Type", "application/json")
	ireq.Header.Set("X-User-ID", followerID.String())

	iResp, err := client.Do(ireq)
	if err != nil {
		return fmt.Errorf("failed to send posts to interaction-service: %w", err)
	}
	defer iResp.Body.Close()

	if iResp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(iResp.Body)
		return fmt.Errorf("interaction-service returned non-200 status: %d. Body: %s", iResp.StatusCode, string(responseBody))
	}

	log.Printf("Successfully populated timeline for user %s with %d posts from %s", followerID, len(timelinePayload), followedID)
	return nil
}

// Follow implements follow.FollowUsecase.
func (f *followUseCase) Follow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)
	if err != nil {
		return err
	}

	userToFollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	if parsedFollowerID == userToFollow.ID {
		return errors.New("you cannot follow yourself")
	}

	if err := f.followRepo.Follow(parsedFollowerID, userToFollow.ID); err != nil {
		return errors.New("user have been followed")
	}

	go func() {
		if err := f.populateTimeline(parsedFollowerID, userToFollow.ID); err != nil {
			log.Printf("ERROR: Failed to populate timeline for follower %s after following %s: %v", parsedFollowerID, userToFollow.ID, err)
		}
	}()

	queueName := "notification.follow.created"
	event := events.FollowCreatedEvent{
		FollowingID: userToFollow.ID,
		SenderID:    parsedFollowerID,
	}

	if err := f.rabbitMQ.PublishJSON(context.Background(), queueName, event); err != nil {
		log.Printf("CRITICAL: Failed to publish follow.created event: %v", err)
	}

	return nil
}

// Unfollow implements follow.FollowUsecase.
func (f *followUseCase) Unfollow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)

	if err != nil {
		return err
	}

	userToUnfollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Unfollow(parsedFollowerID, userToUnfollow.ID)
}

func (f *followUseCase) GetFollowerIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error) {
	return f.followRepo.GetFollowerIDsByUserID(userID)
}
