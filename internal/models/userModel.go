package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username" form:"username"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	CreatedAt string `json:"created_at"`
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email cannot be empty")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}

func (u *User) ComparePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) HashPassword(password string) error {
	bcryptPassword, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if error != nil {
		return error
	}

	u.Password = string(bcryptPassword)
	return nil
}
