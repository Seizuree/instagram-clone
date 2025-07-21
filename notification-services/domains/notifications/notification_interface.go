package notifications

import (
	"notification-services/domains/notifications/entities"
	"notification-services/events"

	"github.com/google/uuid"
)

type NotificationUseCase interface {
	CreateLikeNotification(event *events.LikeCreatedEvent) error
	CreateCommentNotification(event *events.CommentCreatedEvent) error
	CreateFollowNotification(event *events.FollowCreatedEvent) error
	GetNotifications(receiverID uuid.UUID) ([]*entities.Notification, error)
}

type NotificationRepository interface {
	Save(notification *entities.Notification) error
	GetByReceiverID(receiverID uuid.UUID) ([]*entities.Notification, error)
}
