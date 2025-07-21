package entities

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID  `json:"id"`
	ReceiverID uuid.UUID  `json:"receiver_id"`
	SenderID   uuid.UUID  `json:"sender_id"`
	Type       string     `json:"type"` // e.g., "like", "comment", "follow"
	PostID     *uuid.UUID `json:"post_id,omitempty"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
}
