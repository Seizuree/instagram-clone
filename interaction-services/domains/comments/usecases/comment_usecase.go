package usecases

import (
	"context"
	"interaction-services/config"
	"interaction-services/domains/comments"
	"interaction-services/domains/comments/entities"
	"interaction-services/domains/comments/models/responses"
	"interaction-services/events"
	"interaction-services/infrastructures"
	"log"

	"github.com/google/uuid"
)

type commentUseCase struct {
	commentRepo comments.CommentRepository
	rabbitMQ    *infrastructures.RabbitMQ
	config      *config.Config
}

func NewCommentUseCase(commentRepo comments.CommentRepository, rabbitMQ *infrastructures.RabbitMQ, config *config.Config) comments.CommentUseCase {
	return &commentUseCase{commentRepo: commentRepo, rabbitMQ: rabbitMQ, config: config}
}

func (c *commentUseCase) CreateComment(userID, postID uuid.UUID, comment string) (*responses.CommentResponse, error) {
	newComment := &entities.Comment{
		UserID:  userID,
		PostID:  postID,
		Comment: comment,
	}

	createdComment, err := c.commentRepo.CreateComment(newComment)
	if err != nil {
		return nil, err
	}

	queueName := "notification.comment.created"
	event := events.CommentCreatedEvent{
		PostID:   postID,
		SenderID: userID,
		Comment:  comment,
	}
	if err := c.rabbitMQ.PublishJSON(context.Background(), queueName, event); err != nil {
		log.Printf("CRITICAL: Failed to publish comment.created event: %v", err)
	}

	return &responses.CommentResponse{
		ID:      createdComment.ID,
		UserID:  createdComment.UserID,
		Comment: createdComment.Comment,
	}, nil
}

func (c *commentUseCase) GetCommentsByPostID(postID uuid.UUID) ([]*responses.CommentResponse, error) {
	comments, err := c.commentRepo.GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	var commentResponses []*responses.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, &responses.CommentResponse{
			ID:      comment.ID,
			UserID:  comment.UserID,
			Comment: comment.Comment,
		})
	}

	return commentResponses, nil
}

func (c *commentUseCase) CountComments(postID uuid.UUID) (int64, error) {
	return c.commentRepo.CountComments(postID)
}
