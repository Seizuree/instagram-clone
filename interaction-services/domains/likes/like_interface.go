package likes

import "github.com/google/uuid"

type LikeUseCase interface {
	LikePost(userID, postID uuid.UUID) error
	UnlikePost(userID, postID uuid.UUID) error
	CountLikes(postID uuid.UUID) (int64, error)
}

type LikeRepository interface {
	LikePost(userID, postID uuid.UUID) error
	UnlikePost(userID, postID uuid.UUID) error
	CountLikes(postID uuid.UUID) (int64, error)
}
