package wizards

import (
	"interaction-services/config"
	commentHttp "interaction-services/domains/comments/handlers/http"
	commentRepo "interaction-services/domains/comments/repositories"
	commentUc "interaction-services/domains/comments/usecases"
	likeHttp "interaction-services/domains/likes/handlers/http"
	likeRepo "interaction-services/domains/likes/repositories"
	likeUc "interaction-services/domains/likes/usecases"
	"interaction-services/infrastructures"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)
	RabbitMQ         = infrastructures.NewRabbitMQ(Config)

	CommentDatabaseRepo = commentRepo.NewCommentRepository(PostgresDatabase)
	LikeDatabaseRepo    = likeRepo.NewLikeRepository(PostgresDatabase)
	CommentUseCase      = commentUc.NewCommentUseCase(CommentDatabaseRepo)
	LikeUseCase         = likeUc.NewLikeUseCase(LikeDatabaseRepo)
	CommentHttp         = commentHttp.NewCommentHttp(CommentUseCase)
	LikeHttp            = likeHttp.NewLikeHttp(LikeUseCase)
)
