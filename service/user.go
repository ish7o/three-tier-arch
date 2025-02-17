package service

import (
	"fmt"
	"three-tier-arch/models"
	"three-tier-arch/store"
)

type UserService struct {
	userStore *store.UserStore
}

func (s *UserService) DeleteUser(id int) error {
    return s.userStore.DeleteUser(id)
}

func NewUserService(userStore *store.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s *UserService) GetAllUsers() []models.User {
	return s.userStore.GetAllUsers()
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.userStore.GetUser(id)
}

func (s *UserService) CreateUser(input models.UserInput) (*models.User, error) {
	if err := input.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid user input: %w", err)
	}

	return s.userStore.CreateUser(input)
}

func (s *UserService) UpdateUser(id int, input models.UserInput) (*models.User, error) {
	if err := input.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid user input: %w", err)
	}

	return s.userStore.UpdateUser(id, input)
}
