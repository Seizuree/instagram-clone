package repositories

import (
	"notification-services/domains/notifications"
	"notification-services/domains/notifications/entities"
	"notification-services/infrastructures"

	"github.com/google/uuid"
)

// notificationRepository now holds a database connection.
type notificationRepository struct {
	db infrastructures.Database
}

// NewNotificationRepository now requires a database instance.
func NewNotificationRepository(db infrastructures.Database) notifications.NotificationRepository {
	return &notificationRepository{db: db}
}

// Save stores a new notification in the database.
func (r *notificationRepository) Save(notification *entities.Notification) error {
	return r.db.GetInstance().Create(notification).Error
}

// GetByReceiverID retrieves all notifications for a specific user from the database.
func (r *notificationRepository) GetByReceiverID(receiverID uuid.UUID) ([]*entities.Notification, error) {
	var notifications []*entities.Notification
	err := r.db.GetInstance().
		Where("receiver_id = ?", receiverID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}
