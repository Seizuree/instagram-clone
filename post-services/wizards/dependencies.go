package wizards

import (
	"post-services/config"
	"post-services/domains/posts/handlers/http"
	"post-services/domains/posts/repositories"
	"post-services/domains/posts/usecases"
	"post-services/infrastructures"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)
	RabbitMQ         = infrastructures.NewRabbitMQ(Config)
	MinioClient      = infrastructures.NewMinioClient(Config)

	PostDatabaseRepo = repositories.NewPostRepository(PostgresDatabase)
	PostUseCase      = usecases.NewPostUseCase(PostDatabaseRepo, MinioClient)
	PostHttp         = http.NewPostHttp(PostUseCase)
)
