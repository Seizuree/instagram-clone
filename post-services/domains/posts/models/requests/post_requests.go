package requests

type PostUpdateRequest struct {
	Caption string `json:"caption" binding:"required"`
}
