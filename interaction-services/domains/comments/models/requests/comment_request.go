package requests

type CommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}
