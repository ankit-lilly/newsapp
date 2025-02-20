package services

import (
	"context"
	"errors"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/repositories"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("email already registered")
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, user *models.User) (int64, error) {
	if err := user.Validate(); err != nil {
		return -1, err
	}
	err := user.HashPassword(user.Password)

	if err != nil {
		return -1, err
	}

	id, err := s.userRepo.InsertUser(ctx, user)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *UserService) UserExists(ctx context.Context, email string) (bool, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, repositories.NoRecordsFound) {
			return false, nil
		}
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return true, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ValidateUserCredentials(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	valid := user.ComparePasswordHash(password)

	if !valid {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
