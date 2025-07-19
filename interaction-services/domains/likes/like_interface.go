package likes

import "github.com/google/uuid"

type LikeUseCase interface {
	LikePost(userID, postID uuid.UUID) error
	UnlikePost(userID, postID uuid.UUID) error
}

type LikeRepository interface {
	LikePost(userID, postID uuid.UUID) error
	UnlikePost(userID, postID uuid.UUID) error
}
