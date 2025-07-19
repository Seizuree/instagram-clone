package responses

type PostCreationResponse struct {
	UserID   string `json:"userID"`
	ImageURL string `json:"imageURL"`
	Caption  string `json:"caption"`
}

type PostDetailResponse struct {
	ID       string `json:"postID"`
	UserID   string `json:"userIO"`
	ImageURL string `json:"imageURL"`
	Caption  string `json:"caption"`
}

type PostUpdateResponse struct {
	ID      string `json:"postID"`
	Caption string `json:"caption"`
}
