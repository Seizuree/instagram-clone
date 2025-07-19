package repositories

import (
	"core-services/domains/users"
	"core-services/domains/users/entities"
	"core-services/infrastructures"

	"github.com/google/uuid"
)

type userRepository struct {
	db infrastructures.Database
}

func NewUserRepository(db infrastructures.Database) users.UserRepository {
	return &userRepository{db: db}
}

// CreateUser implements users.UserRepository.
func (u *userRepository) CreateUser(user *entities.User) error {
	if err := u.db.GetInstance().Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail implements users.UserRepository.
func (u *userRepository) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := u.db.GetInstance().Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername implements users.UserRepository.
func (u *userRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	if err := u.db.GetInstance().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetUserByID(userID uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := u.db.GetInstance().Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateUser(user *entities.User) error {
	if err := u.db.GetInstance().Save(user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser implements users.UserRepository.
func (u *userRepository) DeleteUser(userID uuid.UUID) error {
	if err := u.db.GetInstance().Where("id = ?", userID).Delete(&entities.User{}).Error; err != nil {
		return err
	}
	return nil
}
