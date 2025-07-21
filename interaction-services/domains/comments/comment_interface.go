package comments

import (
	"interaction-services/domains/comments/entities"
	"interaction-services/domains/comments/models/responses"

	"github.com/google/uuid"
)

type CommentUseCase interface {
	CreateComment(userID, postID uuid.UUID, comment string) (*responses.CommentResponse, error)
	GetCommentsByPostID(postID uuid.UUID) ([]*responses.CommentResponse, error)
	CountComments(postID uuid.UUID) (int64, error)
}

type CommentRepository interface {
	CreateComment(comment *entities.Comment) (*entities.Comment, error)
	GetCommentsByPostID(postID uuid.UUID) ([]*entities.Comment, error)
	CountComments(postID uuid.UUID) (int64, error)
}
