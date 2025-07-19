package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ProxyHeadersMiddleware is a middleware that adds essential headers to requests
// that are about to be proxied to other services.
func ProxyHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Attempt to get the userID from the context (set by the auth middleware).
		userID, exists := c.Get("userID")
		if exists {
			// If the userID exists, add it to the request header so the downstream
			// service knows which user is making the request.
			c.Request.Header.Set("X-User-ID", userID.(string))
		} else {
			// This case should ideally not be hit if the auth middleware is working correctly,
			// but it's good practice to log it if it happens.
			log.Println("WARN: ProxyHeadersMiddleware: userID not found in context.")
		}

		// Proceed to the next handler in the chain (which will be the proxy handler).
		c.Next()
	}
}
