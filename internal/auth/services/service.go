package services

import (
	"github.com/ankibahuguna/newsapp/internal/auth/respository"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepo,
	}
}
