package http

import (
	"net/http"
	"notification-services/domains/notifications"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHttp struct {
	uc notifications.NotificationUseCase
}

func NewNotificationHttp(uc notifications.NotificationUseCase) *NotificationHttp {
	return &NotificationHttp{uc: uc}
}

// GetNotifications fetches the notification history for a user.
func (h *NotificationHttp) GetNotifications(c *gin.Context) {
	// The userID should be passed from the API Gateway's auth middleware.
	userIDstr := c.GetHeader("X-User-ID")
	if userIDstr == "" {
		// Return an error if the header is missing entirely
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID header is missing"})
		return
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing X-User-ID header"})
		return
	}

	notifications, err := h.uc.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
