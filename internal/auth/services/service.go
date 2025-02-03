package services

import (
	"fmt"

	"github.com/ankibahuguna/newsapp/internal/auth/respository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		AuthRepository: authRepo,
	}
}


func (a *AuthService) GetUserByEmail(email string) (*repository.User, error) {
	user, err := a.AuthRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}


func (a *AuthService) CreateUser(u repository.User) (*repository.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashedPassword)

	id, err := a.AuthRepository.InsertUser(u)

	if err != nil {
		return nil, err
	}
	user := &repository.User{Name: u.Name, ID: id}
	return user, err
}

func (a *AuthService) LoginUser(email, password string) (*repository.User, error) {

	user, err := a.AuthRepository.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid password", err)
	}

	return user, nil
}
