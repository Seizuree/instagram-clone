package usecases

import (
	"core-services/config"
	"core-services/domains/users"
	"core-services/domains/users/entities"
	"core-services/domains/users/models/requests"
	"core-services/domains/users/models/responses"
	"core-services/shared/util"
	"errors"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo users.UserRepository
	config   *config.Config
}

func NewUserUseCase(userRepo users.UserRepository) users.UserUseCase {
	return &userUseCase{userRepo: userRepo, config: config.GetConfig()}
}

// Register implements users.UserUsecase.
func (u *userUseCase) Register(req *requests.UserRegisterRequest) (*responses.UserRegisterResponse, error) {
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	newUser := &entities.User{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := u.userRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return &responses.UserRegisterResponse{
		Username: newUser.Username,
		Email:    newUser.Email,
		Message:  "User registered successfully",
	}, nil
}

// Login implements users.UserUsecase.
func (u *userUseCase) Login(req *requests.UserLoginRequest) (*responses.UserLoginResponse, error) {
	user, err := u.userRepo.GetUserByEmail(req.Email)

	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := util.CheckPasswordHash(req.Password, user.Password); !err {
		return nil, errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID.String(), u.config.Server.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &responses.UserLoginResponse{Token: token}, nil
}

// GetProfile implements users.UserUseCase.
func (u *userUseCase) GetProfile(username string) (*responses.UserProfileResponse, error) {
	user, err := u.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &responses.UserProfileResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
