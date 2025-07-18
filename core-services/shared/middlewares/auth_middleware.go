package middlewares

import (
	"core-services/config"
	"core-services/shared/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	conf *config.Config
}

func NewAuthMiddleware(conf *config.Config) AuthMiddleware {
	return AuthMiddleware{conf: conf}
}

func (am *AuthMiddleware) UseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenStrings := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStrings == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		claims, err := util.ValidateJWT(tokenStrings, am.conf.Server.JWTSecret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims["user_id"])
		c.Next()
	}
}
