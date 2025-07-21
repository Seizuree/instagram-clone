package wizards

import (
	websocket_handler "notification-services/domains/notifications/handlers/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	// Register WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		websocket_handler.ServeWs(Hub, c, Config)
	})

	// Register REST API routes
	apiGroup := router.Group("/api")
	{
		notifGroup := apiGroup.Group("/notifications")
		{
			notifGroup.GET("/", NotificationHttp.GetNotifications)
		}
	}
}
