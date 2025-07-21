package events

import (
	"github.com/google/uuid"
)

// This package holds all event structs related to notifications.

type LikeCreatedEvent struct {
	PostID   uuid.UUID `json:"post_id"`
	SenderID uuid.UUID `json:"sender_id"`
}

type CommentCreatedEvent struct {
	PostID   uuid.UUID `json:"post_id"`
	SenderID uuid.UUID `json:"sender_id"`
	Comment  string    `json:"comment"`
}

type FollowCreatedEvent struct {
	FollowingID uuid.UUID `json:"following_id"`
	SenderID    uuid.UUID `json:"sender_id"`
}
