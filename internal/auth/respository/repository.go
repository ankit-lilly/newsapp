package repository

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
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

func (a *AuthRepository) InsertUser(user User) (int, error) {
	query := `INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := a.DB.QueryRow(query, user.Name, user.Email, user.Password, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
