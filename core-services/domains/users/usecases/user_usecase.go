package usecases

import (
	"context"
	"core-services/config"
	"core-services/domains/follow"
	"core-services/domains/users"
	"core-services/domains/users/entities"
	"core-services/domains/users/models/requests"
	"core-services/domains/users/models/responses"
	"core-services/infrastructures"
	"core-services/shared/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo   users.UserRepository
	followRepo follow.FollowRepository
	config     *config.Config
	rabbitMQ   *infrastructures.RabbitMQ
}

func NewUserUseCase(userRepo users.UserRepository, followRepo follow.FollowRepository, rabbitMQ *infrastructures.RabbitMQ) users.UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		followRepo: followRepo,
		config:     config.GetConfig(),
		rabbitMQ:   rabbitMQ,
	}
}

// Register implements users.UserUsecase.
func (u *userUseCase) Register(req *requests.UserRegisterRequest) (*responses.UserRegisterResponse, error) {
	_, err := u.userRepo.GetUserByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	// Check if email is already taken
	_, err = u.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already taken")
	}

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

	followerCount, followingCount, err := u.followRepo.GetFollowCounts(user.ID)

	if err != nil {
		return nil, err
	}

	postCount, err := u.getPostCount(user.ID)
	if err != nil {
		// You might want to decide if this should be a fatal error or if you
		// can return a response with a post count of 0.
		log.Printf("Failed to get post count for user %s: %v", user.ID, err)
	}

	return &responses.UserProfileResponse{
		Username:       user.Username,
		Email:          user.Email,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		PostCount:      postCount,
	}, nil
}

func (u *userUseCase) getPostCount(userID uuid.UUID) (int64, error) {
	postServiceURL := u.config.Server.PostServiceURL
	resp, err := http.Get(fmt.Sprintf("%s/api/posts/user/%s/count", postServiceURL, userID.String()))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("post-service returned non-OK status: %d", resp.StatusCode)
	}

	var result struct {
		Count int64 `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Count, nil
}

func (u *userUseCase) GetMe(userID uuid.UUID) (*responses.UserProfileResponse, error) {
	user, err := u.userRepo.GetUserByID(userID)

	if err != nil {
		return nil, errors.New("user not found")
	}

	followerCount, followingCount, err := u.followRepo.GetFollowCounts(user.ID)

	if err != nil {
		return nil, err
	}

	postCount, err := u.getPostCount(user.ID)
	if err != nil {
		// You might want to decide if this should be a fatal error or if you
		// can return a response with a post count of 0.
		log.Printf("Failed to get post count for user %s: %v", user.ID, err)
	}

	return &responses.UserProfileResponse{
		Username:       user.Username,
		Email:          user.Email,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		PostCount:      postCount,
	}, nil
}

// UpdateUser implements users.UserUseCase.
func (u *userUseCase) UpdateUser(userID uuid.UUID, req *requests.UserUpdateRequest) (*responses.UserUpdateResponse, error) {
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := u.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &responses.UserUpdateResponse{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// DeleteUser implements users.UserUseCase.
func (u *userUseCase) DeleteUser(userID uuid.UUID) error {
	if err := u.userRepo.DeleteUser(userID); err != nil {
		log.Printf("Error deleting user from repository: %v", err)
		return err
	}

	// Publish a message to RabbitMQ that a user has been deleted.
	// This is an asynchronous operation.
	// The post-service will listen for this message to delete related posts.
	queueName := "user.deleted"
	message := map[string]interface{}{"user_id": userID.String()}

	if err := u.rabbitMQ.PublishJSON(context.Background(), queueName, message); err != nil {
		// Log the error but don't return it to the user.
		// The core operation (user deletion) was successful.
		// We need a more robust out-of-band error handling mechanism for this (e.g., monitoring, alerts).
		log.Printf("CRITICAL: Failed to publish user.deleted event for userID %s: %v", userID, err)
	}

	log.Printf("Successfully deleted user %s and published event to %s", userID, queueName)
	return nil
}
