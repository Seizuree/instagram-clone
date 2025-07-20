package notifications

import (
	"notification-service/entities"

	"github.com/google/uuid"
)

type NotificationUseCase interface {
	SendNotification(notification *entities.Notification) error
	GetNotifications(receiverID uuid.UUID) ([]*entities.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
}

type NotificationRepository interface {
	Save(notification *entities.Notification) error
	GetByReceiverID(receiverID uuid.UUID) ([]*entities.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
}
