package wizards

import (
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		posts := apiGroup.Group("/posts")
		{
			posts.POST("", PostHttp.CreatePost)
			posts.GET("/:post_id", PostHttp.GetPost)
			posts.GET("/user/:user_id", PostHttp.GetPostsByUser)
			posts.PUT("/:post_id", PostHttp.UpdatePost)
			posts.DELETE("/:post_id", PostHttp.DeletePost)
		}
	}
}
