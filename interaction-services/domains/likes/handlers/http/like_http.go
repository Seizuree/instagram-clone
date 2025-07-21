package http

import (
	"errors"
	"interaction-services/domains/likes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeHttp struct {
	lc likes.LikeUseCase
}

func NewLikeHttp(lc likes.LikeUseCase) *LikeHttp {
	return &LikeHttp{lc}
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		return uuid.Nil, errors.New("X-User-ID header is missing, request must come from the API gateway")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user id format in X-User-ID header")
	}
	return userID, nil
}

func (h *LikeHttp) LikePost(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.lc.LikePost(userID, postID); err != nil {
		if strings.Contains(err.Error(), "post not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already liked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post liked successfully"})
}

func (h *LikeHttp) UnlikePost(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.lc.UnlikePost(userID, postID); err != nil {
		if strings.Contains(err.Error(), "not liked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post unliked successfully"})
}
