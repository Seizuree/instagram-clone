package http

import (
	"core-services/domains/follow"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowHttp struct {
	fc follow.FollowUseCase
}

func NewFollowHttp(fc follow.FollowUseCase) *FollowHttp {
	return &FollowHttp{fc}
}

func (h *FollowHttp) FollowUser(c *gin.Context) {
	username := c.Param("username")
	userID, _ := c.Get("user_id")

	if err := h.fc.Follow(userID.(uuid.UUID), username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user followed successfully"})
}

func (h *FollowHttp) UnfollowUser(c *gin.Context) {
	username := c.Param("username")
	userID, _ := c.Get("user_id")

	if err := h.fc.Unfollow(userID.(uuid.UUID), username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unfollowed successfully"})
}
