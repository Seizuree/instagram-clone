package wizards

import (
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		interactions := apiGroup.Group("/interactions")
		{
			// Like and Unlike a post
			interactions.POST("/:post_id/like", LikeHttp.LikePost)
			interactions.DELETE("/:post_id/like", LikeHttp.UnlikePost)

			// Comment on a post
			interactions.POST("/:post_id/comment", CommentHttp.CreateComment)
			interactions.GET("/:post_id/comments", CommentHttp.GetCommentsByPostID)
			
			// folow
			interactions.POST("/follow", followHtttp.Follow)
			interactions.DELETE("/unfollow", followHttp.Unfollow)

		}
	}
}
