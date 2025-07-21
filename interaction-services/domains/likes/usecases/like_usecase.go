package usecases

import (
	"context"
	"interaction-services/config"
	"interaction-services/domains/likes"
	"interaction-services/events"
	"interaction-services/infrastructures"
	"log"

	"github.com/google/uuid"
)

type likeUseCase struct {
	likeRepo likes.LikeRepository
	rabbitMQ *infrastructures.RabbitMQ
	config   *config.Config
}

func NewLikeUseCase(likeRepo likes.LikeRepository, rabbitMQ *infrastructures.RabbitMQ, config *config.Config) likes.LikeUseCase {
	return &likeUseCase{likeRepo: likeRepo, rabbitMQ: rabbitMQ, config: config}
}

func (l *likeUseCase) LikePost(userID, postID uuid.UUID) error {
	if err := l.likeRepo.LikePost(userID, postID); err != nil {
		return err
	}

	queueName := "notification.like.created"
	event := events.LikeCreatedEvent{
		PostID:   postID,
		SenderID: userID,
	}

	if err := l.rabbitMQ.PublishJSON(context.Background(), queueName, event); err != nil {
		log.Printf("CRITICAL: Failed to publish like.created event: %v", err)
	}

	return nil
}

func (l *likeUseCase) UnlikePost(userID, postID uuid.UUID) error {
	return l.likeRepo.UnlikePost(userID, postID)
}

func (l *likeUseCase) CountLikes(postID uuid.UUID) (int64, error) {
	return l.likeRepo.CountLikes(postID)
}
