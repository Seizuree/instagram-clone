package http

import (
	"interaction-services/domains/comments"
	"interaction-services/domains/likes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InternalHttp struct {
	likeUseCase    likes.LikeUseCase
	commentUseCase comments.CommentUseCase
}

func NewInternalHttp(likeUc likes.LikeUseCase, commentUc comments.CommentUseCase) *InternalHttp {
	return &InternalHttp{
		likeUseCase:    likeUc,
		commentUseCase: commentUc,
	}
}

func (h *InternalHttp) GetInteractionCounts(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	likeCount, err := h.likeUseCase.CountLikes(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count likes"})
		return
	}

	commentCount, err := h.commentUseCase.CountComments(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"like_count":    likeCount,
		"comment_count": commentCount,
	})
}
