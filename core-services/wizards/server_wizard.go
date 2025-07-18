package wizards

import (
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/register", UserHttp.Register)
		apiGroup.POST("/login", UserHttp.Login)
		apiGroup.Use(AuthMiddleware.UseAuthMiddleware())
		user := apiGroup.Group("/users")
		{
			user.GET("/:username", UserHttp.GetProfile)
			user.POST("/follow/:username", FollowHttp.FollowUser)
			user.POST("/unfollow/:username", FollowHttp.UnfollowUser)
		}
	}
}
