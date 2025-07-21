package http

import (
	"interaction-services/domains/timelines"
	"interaction-services/domains/timelines/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TimelineHttp struct {
	usecase timelines.TimelineUseCase
}

func NewTimelineHttp(u timelines.TimelineUseCase) *TimelineHttp {
	return &TimelineHttp{usecase: u}
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, error) {
	userIDStr := c.GetHeader("X-User-ID")
	return uuid.Parse(userIDStr)
}

func (h *TimelineHttp) GetTimeline(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing X-User-ID"})
		return
	}

	timeline, err := h.usecase.GetTimeline(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get timeline"})
		return
	}

	c.JSON(http.StatusOK, timeline)
}

func (h *TimelineHttp) AddPostsToTimeline(c *gin.Context) {
	userID, err := getUserIDFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing X-User-ID"})
		return
	}

	var posts []*entities.Timeline
	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.usecase.AddPostsToFollowerTimeline(userID, posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add posts to timeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "timeline updated successfully"})
}
