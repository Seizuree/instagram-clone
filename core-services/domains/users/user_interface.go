package users

import (
	"core-services/domains/users/entities"
	"core-services/domains/users/models/requests"
	"core-services/domains/users/models/responses"
)

type UserUseCase interface {
	Register(req *requests.UserRegisterRequest) (*responses.UserRegisterResponse, error)
	Login(req *requests.UserLoginRequest) (*responses.UserLoginResponse, error)
	GetProfile(username string) (*responses.UserProfileResponse, error)
}

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
}
