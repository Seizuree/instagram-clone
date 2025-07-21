package wizards

import (
	"core-services/config"
	followHttp "core-services/domains/follow/handlers/http"
	followRepo "core-services/domains/follow/repositories"
	followUc "core-services/domains/follow/usecases"
	userHttp "core-services/domains/users/handlers/http"
	userRepo "core-services/domains/users/repositories"
	userUc "core-services/domains/users/usecases"
	"core-services/infrastructures"
	"core-services/shared/middlewares"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)
	RabbitMQ         = infrastructures.NewRabbitMQ(Config)

	UserDatabaseRepo   = userRepo.NewUserRepository(PostgresDatabase)
	FollowDatabaseRepo = followRepo.NewFollowRepository(PostgresDatabase)
	UserUseCase        = userUc.NewUserUseCase(UserDatabaseRepo, FollowDatabaseRepo, RabbitMQ)
	FollowUseCase      = followUc.NewFollowUseCase(FollowDatabaseRepo, UserDatabaseRepo, RabbitMQ, Config)
	UserHttp           = userHttp.NewUserHttp(UserUseCase)
	FollowHttp         = followHttp.NewFollowHttp(FollowUseCase)

	AuthMiddleware  = middlewares.NewAuthMiddleware(Config)
	ProxyMiddleware = middlewares.ProxyHeadersMiddleware()
)
