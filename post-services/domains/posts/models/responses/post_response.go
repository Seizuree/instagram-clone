package responses

type PostCreationResponse struct {
	UserID   string `json:"userID"`
	ImageURL string `json:"imageURL"`
	Caption  string `json:"caption"`
}

type PostDetailResponse struct {
	ID           string `json:"postID"`
	UserID       string `json:"userID"`
	ImageURL     string `json:"imageURL"`
	Caption      string `json:"caption"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
}

type PostUpdateResponse struct {
	ID      string `json:"postID"`
	Caption string `json:"caption"`
}

type PostInteractionCounts struct {
	LikeCount    int64 `json:"like_count"`
	CommentCount int64 `json:"comment_count"`
}
