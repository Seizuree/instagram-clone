package wizards

import (
	"interaction-services/domains/comments/http"
	CommentHttp "interaction-services/domains/comments/http"
	LikeHttp "interaction-services/domains/likes/http"
	FollowHttp "interaction-services/domains/follows/http"
	TimelineHttp "interaction-services/domains/timelines/http"

	CommentRepo "interaction-services/domains/comments/repositories"
	LikeRepo "interaction-services/domains/likes/repositories"
	FollowRepo "interaction-services/domains/follows/repositories"
	TimelineRepo "interaction-services/domains/timelines/repositories"

	CommentUsecase "interaction-services/domains/comments/usecases"
	LikeUsecase "interaction-services/domains/likes/usecases"
	FollowUsecase "interaction-services/domains/follows/usecases"
	TimelineUsecase "interaction-services/domains/timelines/usecases"

	"interaction-services/infrastructures"

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine, db infrastructures.Database) {
	apiGroup := router.Group("/api")

	// === INTERACTIONS ===
	interactions := apiGroup.Group("/interactions")
	{
		// LIKE
		likeRepo := LikeRepo.NewLikeRepository(db)
		likeUC := LikeUsecase.NewLikeUseCase(likeRepo)
		likeHandler := LikeHttp.NewLikeHandler(likeUC)
		interactions.POST("/:post_id/like", likeHandler.LikePost)
		interactions.DELETE("/:post_id/like", likeHandler.UnlikePost)

		// COMMENT
		commentRepo := CommentRepo.NewCommentRepository(db)
		commentUC := CommentUsecase.NewCommentUseCase(commentRepo)
		commentHandler := CommentHttp.NewCommentHandler(commentUC)
		interactions.POST("/:post_id/comment", commentHandler.CreateComment)
		interactions.GET("/:post_id/comments", commentHandler.GetCommentsByPostID)

		// FOLLOW
		followRepo := FollowRepo.NewFollowRepository(db)
		followUC := FollowUsecase.NewFollowUseCase(followRepo)
		followHandler := FollowHttp.NewFollowHandler(followUC)
		interactions.POST("/follow", followHandler.Follow)
		interactions.DELETE("/unfollow", followHandler.Unfollow)
	}

	// === TIMELINE ===
	timelineRepo := TimelineRepo.NewTimelineRepository(db)
	timelineUC := TimelineUsecase.NewTimelineUseCase(timelineRepo)
	timelineHandler := TimelineHttp.NewTimelineHandler(timelineUC)

	timeline := apiGroup.Group("/timeline")
	{
		timeline.GET("/", timelineHandler.GetTimeline)
	}
}
