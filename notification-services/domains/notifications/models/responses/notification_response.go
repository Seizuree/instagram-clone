package responses

import "github.com/google/uuid"

type PostResponse struct {
	UserID uuid.UUID `json:"userID"`
}

type UserResponse struct {
	Username string `json:"username"`
}
