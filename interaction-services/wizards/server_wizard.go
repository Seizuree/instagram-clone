package wizards

import "github.com/gin-gonic/gin"

func RegisterServer(router *gin.Engine) {
	apiGroup := router.Group("/api")

	// === INTERACTIONS ===
	interactions := apiGroup.Group("/interactions")
	{
		// LIKE
		interactions.POST("/:post_id/like", LikeHttp.LikePost)
		interactions.DELETE("/:post_id/like", LikeHttp.UnlikePost)

		// COMMENT
		interactions.POST("/:post_id/comment", CommentHttp.CreateComment)
		interactions.GET("/:post_id/comments", CommentHttp.GetCommentsByPostID)
		// TIMELINE - Move this route here
		interactions.GET("/timeline", TimelineHttp.GetTimeline)
	}

	internalGroup := router.Group("/api/internal")
	{
		internalGroup.GET("/interactions/:post_id/counts", InternalHttp.GetInteractionCounts)
	}
}
