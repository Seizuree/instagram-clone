package wizards

import (
	"notification-services/config"
	"notification-services/domains/notifications/handlers/http"
	"notification-services/domains/notifications/handlers/websocket"
	"notification-services/domains/notifications/repositories"
	"notification-services/domains/notifications/usecases"
	"notification-services/infrastructures"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)
	RabbitMQ         = infrastructures.NewRabbitMQ(Config)
	Hub              = websocket.NewHub()

	NotificationDatabaseRepo = repositories.NewNotificationRepository(PostgresDatabase)
	NotificationUseCase      = usecases.NewNotificationUseCase(NotificationDatabaseRepo, Hub, Config)
	NotificationHttp         = http.NewNotificationHttp(NotificationUseCase)
)
