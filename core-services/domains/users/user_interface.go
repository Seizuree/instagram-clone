package users

import (
	"core-services/domains/users/entities"
	"core-services/domains/users/models/requests"
	"core-services/domains/users/models/responses"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Register(req *requests.UserRegisterRequest) (*responses.UserRegisterResponse, error)
	Login(req *requests.UserLoginRequest) (*responses.UserLoginResponse, error)
	GetMe(userID uuid.UUID) (*responses.UserProfileResponse, error)
	GetProfile(username string) (*responses.UserProfileResponse, error)
	UpdateUser(userID uuid.UUID, req *requests.UserUpdateRequest) (*responses.UserProfileResponse, error)
	DeleteUser(userID uuid.UUID) error
}

type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByID(userID uuid.UUID) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateUser(user *entities.User) error
	DeleteUser(userID uuid.UUID) error
}
