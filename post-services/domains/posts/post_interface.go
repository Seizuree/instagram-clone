package posts

import (
	"mime/multipart"
	"post-services/domains/posts/entities"
	"post-services/domains/posts/models/responses"

	"github.com/google/uuid"
)

type PostUseCase interface {
	CreatePost(userID uuid.UUID, caption string, fileHeader *multipart.FileHeader) (*entities.Post, error)
	GetPost(postID uuid.UUID) (*responses.PostDetailResponse, error)
	GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error)
	UpdatePost(userID, postID uuid.UUID, caption string) (*entities.Post, error)
	DeletePost(userID, postID uuid.UUID) error
	DeletePostsByUserID(userID uuid.UUID) error
	CountUserPosts(userID uuid.UUID) (int64, error)
}

type PostRepository interface {
	CreatePost(post *entities.Post) error
	GetPostByID(postID uuid.UUID) (*entities.Post, error)
	GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error)
	UpdatePost(post *entities.Post) error
	DeletePost(postID uuid.UUID) error
	DeletePostsByUserID(userID uuid.UUID) error
	CountUserPosts(userID uuid.UUID) (int64, error)
}
