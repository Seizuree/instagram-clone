package usecases

import (
	"interaction-services/domains/likes"

	"github.com/google/uuid"
)

type likeUseCase struct {
	likeRepo likes.LikeRepository
}

func NewLikeUseCase(likeRepo likes.LikeRepository) likes.LikeUseCase {
	return &likeUseCase{likeRepo: likeRepo}
}

func (l *likeUseCase) LikePost(userID, postID uuid.UUID) error {
	return l.likeRepo.LikePost(userID, postID)
}

func (l *likeUseCase) UnlikePost(userID, postID uuid.UUID) error {
	return l.likeRepo.UnlikePost(userID, postID)
}
