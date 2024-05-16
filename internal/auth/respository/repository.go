package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AuthRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (a *AuthRepository) InsertUser(user User) (int64, error) {
	query := `INSERT INTO users (name, email, password ) VALUES ($1, $2, $3)`
	res, err := a.DB.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := res.LastInsertId()
	return lastInsertId, err
}

func (a *AuthRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`

	var user User
	err := a.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with email %s not found", email)
		}
		return nil, err
	}
	return &user, nil
}
