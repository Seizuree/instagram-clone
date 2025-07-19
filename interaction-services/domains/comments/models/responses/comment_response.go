package responses

import "github.com/google/uuid"

type CommentResponse struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userID"`
	Comment string    `json:"comment"`
}
