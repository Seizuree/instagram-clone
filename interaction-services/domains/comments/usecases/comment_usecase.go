package usecases

import (
	"interaction-services/domains/comments"
	"interaction-services/domains/comments/entities"
	"interaction-services/domains/comments/models/responses"

	"github.com/google/uuid"
)

type commentUseCase struct {
	commentRepo comments.CommentRepository
}

func NewCommentUseCase(commentRepo comments.CommentRepository) comments.CommentUseCase {
	return &commentUseCase{commentRepo: commentRepo}
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
