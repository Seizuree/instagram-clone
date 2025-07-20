package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"interaction-services/domains/follows/usecases"
)

type FollowHandler struct {
	usecase *usecases.FollowUsecase
}

func NewFollowHandler(usecase *usecases.FollowUsecase) *FollowHandler {
	return &FollowHandler{usecase: usecase}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	followerID := c.Query("follower_id")
	followingID := c.Query("following_id")

	followerUUID, _ := uuid.Parse(followerID)
	followingUUID, _ := uuid.Parse(followingID)

	err := h.usecase.Follow(followerUUID, followingUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID := c.Query("follower_id")
	followingID := c.Query("following_id")

	followerUUID, _ := uuid.Parse(followerID)
	followingUUID, _ := uuid.Parse(followingID)

	err := h.usecase.Unfollow(followerUUID, followingUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}
