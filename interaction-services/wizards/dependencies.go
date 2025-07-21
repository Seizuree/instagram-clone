package wizards

import (
	"interaction-services/config"
	commentHttp "interaction-services/domains/comments/handlers/http"
	commentRepo "interaction-services/domains/comments/repositories"
	commentUc "interaction-services/domains/comments/usecases"
	internalHttp "interaction-services/domains/interactions/handlers/http"
	likeHttp "interaction-services/domains/likes/handlers/http"
	likeRepo "interaction-services/domains/likes/repositories"
	likeUc "interaction-services/domains/likes/usecases"
	timelineHttp "interaction-services/domains/timelines/handlers/http"
	timelineRepo "interaction-services/domains/timelines/repositories"
	timelineUc "interaction-services/domains/timelines/usecases"
	"interaction-services/infrastructures"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)
	RabbitMQ         = infrastructures.NewRabbitMQ(Config)

	CommentDatabaseRepo  = commentRepo.NewCommentRepository(PostgresDatabase)
	LikeDatabaseRepo     = likeRepo.NewLikeRepository(PostgresDatabase)
	TimelineDatabaseRepo = timelineRepo.NewTimelineRepository(PostgresDatabase)

	CommentUseCase  = commentUc.NewCommentUseCase(CommentDatabaseRepo, RabbitMQ, Config)
	LikeUseCase     = likeUc.NewLikeUseCase(LikeDatabaseRepo, RabbitMQ, Config)
	TimelineUseCase = timelineUc.NewTimelineUseCase(TimelineDatabaseRepo, Config)

	CommentHttp  = commentHttp.NewCommentHttp(CommentUseCase)
	LikeHttp     = likeHttp.NewLikeHttp(LikeUseCase)
	TimelineHttp = timelineHttp.NewTimelineHttp(TimelineUseCase)
	InternalHttp = internalHttp.NewInternalHttp(LikeUseCase, CommentUseCase)
)
