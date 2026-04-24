package service

import (
	"context"
	"errors"
	"fmt"
	"ums/model"
	"ums/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) RegisterUser(ctx context.Context, name, email string) (*model.User, error) {
	allUsers, err := s.repo.List(ctx)

	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}

	for _, u := range allUsers {
		if u.Email == email {
			return nil, errors.New("email already in use")
		}
	}

	user := &model.User{Name: name, Email: email}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, name, email string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Name = name
	user.Email = email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
