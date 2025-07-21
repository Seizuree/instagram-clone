package http

import (
	"core-services/domains/follow"
	"net/http"
	"strings"

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
	userID, _ := c.Get("userID")

	if err := h.fc.Follow(userID.(string), username); err != nil {
		if strings.Contains(err.Error(), "cannot follow yourself") || strings.Contains(err.Error(), "have been") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user followed successfully"})
}

func (h *FollowHttp) UnfollowUser(c *gin.Context) {
	username := c.Param("username")
	userID, _ := c.Get("userID")

	if err := h.fc.Unfollow(userID.(string), username); err != nil {
		if strings.Contains(err.Error(), "already unfollowed") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Handle other potential errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unfollowed successfully"})
}

func (h *FollowHttp) GetFollowers(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	followerIDs, err := h.fc.GetFollowerIDsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"follower_ids": followerIDs})
}
