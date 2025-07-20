package wizards

import (
	"github.com/gin-gonic/gin"
	"post-services/http" // adjust to actual path if needed
)

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		postGroup := apiGroup.Group("/posts")
		{
			postGroup.POST("", http.PostHttp.CreatePost)                  // Upload image & generate thumbnail
			postGroup.GET("/:post_id", http.PostHttp.GetPost)            // Get single post
			postGroup.GET("/user/:user_id", http.PostHttp.GetPostsByUser) // Get all posts by user
			postGroup.PUT("/:post_id", http.PostHttp.UpdatePost)         // Update caption
			postGroup.DELETE("/:post_id", http.PostHttp.DeletePost)      // Delete post

			// 🔧 New planned features
			postGroup.GET("/timeline", http.PostHttp.GetTimeline)        // Timeline generation
			postGroup.GET("/user/:user_id/count", http.PostHttp.CountUserPosts) // Post count
		}
	}
}
