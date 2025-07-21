package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notification-services/config"
	"notification-services/domains/notifications"
	"notification-services/domains/notifications/entities"
	"notification-services/domains/notifications/handlers/websocket"
	"notification-services/domains/notifications/models/responses"
	"notification-services/events"
	"time"

	"github.com/google/uuid"
)

type notificationUseCase struct {
	repo notifications.NotificationRepository
	hub  *websocket.Hub
	cfg  *config.Config
}

func NewNotificationUseCase(
	repo notifications.NotificationRepository,
	hub *websocket.Hub,
	cfg *config.Config,
) notifications.NotificationUseCase {
	return &notificationUseCase{
		repo: repo,
		hub:  hub,
		cfg:  cfg,
	}
}

// getPostOwnerID makes a real API call to post-services to find the owner of a post.
func (uc *notificationUseCase) getPostOwnerID(userID, postID uuid.UUID) (uuid.UUID, error) {
	endpoint := fmt.Sprintf("%s/api/internal/posts/%s", uc.cfg.Server.PostServiceURL, postID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-User-ID", userID.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to contact post-service: %w", err)
	}

	defer resp.Body.Close()

	var postResponse responses.PostResponse
	if err := json.NewDecoder(resp.Body).Decode(&postResponse); err != nil {
		return uuid.Nil, fmt.Errorf("failed to unmarshal post response: %w", err)
	}

	return postResponse.UserID, nil
}

// getUsername makes a real API call to core-services to get a user's username.
func (uc *notificationUseCase) getUsername(userID uuid.UUID) (string, error) {
	endpoint := fmt.Sprintf("%s/api/internal/users/%s", uc.cfg.Server.CoreServiceURL, userID)
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to call core-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("core-service returned non-200 status: %d", resp.StatusCode)
	}

	var userResponse responses.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal user response: %w", err)
	}

	return userResponse.Username, nil
}

// CreateLikeNotification handles the logic for a "like" event.
func (uc *notificationUseCase) CreateLikeNotification(event *events.LikeCreatedEvent) error {
	receiverID, err := uc.getPostOwnerID(event.SenderID, event.PostID)
	if err != nil {
		return fmt.Errorf("could not get post owner for like event: %w", err)
	}

	if receiverID == event.SenderID {
		return nil
	}

	senderUsername, err := uc.getUsername(event.SenderID)
	if err != nil {
		log.Printf("WARN: Could not get username for sender %s: %v", event.SenderID, err)
		senderUsername = "Someone" // Fallback username
	}

	notification := &entities.Notification{
		ID:         uuid.New(),
		ReceiverID: receiverID,
		SenderID:   event.SenderID,
		Type:       "like",
		PostID:     &event.PostID,
		Message:    fmt.Sprintf("%s liked your post.", senderUsername),
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.Save(notification); err != nil {
		return err
	}
	uc.hub.PushToUser(receiverID, notification)
	return nil
}

// CreateCommentNotification handles the logic for a "comment" event.
func (uc *notificationUseCase) CreateCommentNotification(event *events.CommentCreatedEvent) error {
	receiverID, err := uc.getPostOwnerID(event.SenderID, event.PostID)
	if err != nil {
		return fmt.Errorf("could not get post owner for comment event: %w", err)
	}

	if receiverID == event.SenderID {
		return nil
	}

	senderUsername, err := uc.getUsername(event.SenderID)
	if err != nil {
		log.Printf("WARN: Could not get username for sender %s: %v", event.SenderID, err)
		senderUsername = "Someone"
	}

	notification := &entities.Notification{
		ID:         uuid.New(),
		ReceiverID: receiverID,
		SenderID:   event.SenderID,
		Type:       "comment",
		PostID:     &event.PostID,
		Message:    fmt.Sprintf(`%s commented: "%s"`, senderUsername, event.Comment),
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.Save(notification); err != nil {
		return err
	}
	uc.hub.PushToUser(receiverID, notification)
	return nil
}

// CreateFollowNotification handles the logic for a "follow" event.
func (uc *notificationUseCase) CreateFollowNotification(event *events.FollowCreatedEvent) error {
	senderUsername, err := uc.getUsername(event.SenderID)
	if err != nil {
		log.Printf("WARN: Could not get username for sender %s: %v", event.SenderID, err)
		senderUsername = "Someone"
	}

	notification := &entities.Notification{
		ID:         uuid.New(),
		ReceiverID: event.FollowingID,
		SenderID:   event.SenderID,
		Type:       "follow",
		PostID:     nil,
		Message:    fmt.Sprintf("%s started following you.", senderUsername),
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.Save(notification); err != nil {
		return err
	}
	uc.hub.PushToUser(event.FollowingID, notification)
	return nil
}

func (uc *notificationUseCase) GetNotifications(receiverID uuid.UUID) ([]*entities.Notification, error) {
	return uc.repo.GetByReceiverID(receiverID)
}
