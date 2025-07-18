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
	Username string `json:"username"`
	Email    string `json:"email"`
}
