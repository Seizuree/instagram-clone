package wizards

import "github.com/gin-gonic/gin"

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		postGroup := apiGroup.Group("/posts")
		{
			postGroup.POST("/", PostHttp.CreatePost)
			postGroup.GET("/:post_id", PostHttp.GetPost)
			postGroup.GET("/user/:user_id", PostHttp.GetPostsByUser)
			postGroup.PUT("/:post_id", PostHttp.UpdatePost)
			postGroup.DELETE("/:post_id", PostHttp.DeletePost)
			postGroup.GET("/user/:user_id/count", PostHttp.CountUserPosts)
		}
	}

	// Add the new internal route group for service-to-service communication
	internalGroup := router.Group("/api/internal")
	{
		internalGroup.GET("/posts/:post_id", PostHttp.GetPost)
	}
}
