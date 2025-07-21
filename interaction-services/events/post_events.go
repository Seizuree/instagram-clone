package events

import (
	"time"

	"github.com/google/uuid"
)

// PostCreatedEvent represents the message structure for a new post.
// It lives in its own package to be safely imported by any other service package.
type PostCreatedEvent struct {
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	ImageURL  string    `json:"image_url"`
	ThumbURL  string    `json:"thumb_url"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeCreatedEvent struct {
	PostID   uuid.UUID `json:"post_id"`
	SenderID uuid.UUID `json:"sender_id"` // Same as UserID, for consistency
}

type CommentCreatedEvent struct {
	PostID   uuid.UUID `json:"post_id"`
	SenderID uuid.UUID `json:"sender_id"` // Same as UserID
	Comment  string    `json:"comment"`
}
