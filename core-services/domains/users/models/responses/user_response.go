package responses

type UserRegisterResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserProfileResponse struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"following_count"`
	PostCount      int64  `json:"post_count"`
}

type UserUpdateResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserPostCountResponse struct {
	Count int64 `json:"count"`
}
