package posts

import (
	"mime/multipart"
	"post-services/domains/posts/entities"

	"github.com/google/uuid"
)

type PostUseCase interface {
	CreatePost(userID uuid.UUID, caption string, fileHeader *multipart.FileHeader) (*entities.Post, error)
	GetPost(postID uuid.UUID) (*entities.Post, error)
	GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error)
	UpdatePost(userID, postID uuid.UUID, caption string) (*entities.Post, error)
	DeletePost(userID, postID uuid.UUID) error
	DeletePostsByUserID(userID uuid.UUID) error
	Save(post *Post) error
	CountByUser(userID uuid.UUID) (int64, error)

	// 🆕 New features
	GetTimeline(userID uuid.UUID) ([]entities.Post, error) // For timeline feature
}

type PostRepository interface {
	CreatePost(post *entities.Post) error
	GetPostByID(postID uuid.UUID) (*entities.Post, error)
	GetPostsByUserID(userID uuid.UUID) (*[]entities.Post, error)
	UpdatePost(post *entities.Post) error
	DeletePost(postID uuid.UUID) error
	DeletePostsByUserID(userID uuid.UUID) error
	Save(post *Post) error
	CountByUser(userID uuid.UUID) (int64, error)

	// 🆕 New feature
	GetTimeline(userID uuid.UUID) ([]entities.Post, error)
}
