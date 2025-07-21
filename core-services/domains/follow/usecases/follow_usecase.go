package usecases

import (
	"context"
	"core-services/domains/follow"
	"core-services/domains/users"
	"core-services/events"
	"core-services/infrastructures"
	"errors"
	"log"

	"github.com/google/uuid"
)

type followUseCase struct {
	followRepo follow.FollowRepository
	userRepo   users.UserRepository
	rabbitMQ   *infrastructures.RabbitMQ
}

func NewFollowUseCase(followRepo follow.FollowRepository, userRepo users.UserRepository, rabbitMQ *infrastructures.RabbitMQ) follow.FollowUseCase {
	return &followUseCase{followRepo: followRepo, userRepo: userRepo, rabbitMQ: rabbitMQ}
}

// Follow implements follow.FollowUsecase.
func (f *followUseCase) Follow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)
	if err != nil {
		return err
	}

	userToFollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}

	if parsedFollowerID == userToFollow.ID {
		return errors.New("you cannot follow yourself")
	}

	if err := f.followRepo.Follow(parsedFollowerID, userToFollow.ID); err != nil {
		return errors.New("user have been followed")
	}

	queueName := "notification.follow.created"
	event := events.FollowCreatedEvent{
		FollowingID: userToFollow.ID,
		SenderID:    parsedFollowerID,
	}

	if err := f.rabbitMQ.PublishJSON(context.Background(), queueName, event); err != nil {
		log.Printf("CRITICAL: Failed to publish follow.created event: %v", err)
	}

	return nil
}

// Unfollow implements follow.FollowUsecase.
func (f *followUseCase) Unfollow(followerID string, username string) error {
	parsedFollowerID, err := uuid.Parse(followerID)

	if err != nil {
		return err
	}

	userToUnfollow, err := f.userRepo.GetUserByUsername(username)
	if err != nil {
		return err
	}
	return f.followRepo.Unfollow(parsedFollowerID, userToUnfollow.ID)
}

func (f *followUseCase) GetFollowerIDsByUserID(userID uuid.UUID) ([]uuid.UUID, error) {
	return f.followRepo.GetFollowerIDsByUserID(userID)
}
