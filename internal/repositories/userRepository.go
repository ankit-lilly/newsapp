package repositories

import (
	"context"
	"database/sql"
	"github.com/ankit-lilly/newsapp/internal/models"
	"log/slog"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{DB: db}
}

func (a *userRepository) InsertUser(ctx context.Context, user *models.User) (int64, error) {
	query := `INSERT INTO users (username, email, password ) VALUES ($1, $2, $3)`
	res, err := a.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		slog.Error("Error inserting user", "error", err)
		return 0, err
	}
	lastInsertId, err := res.LastInsertId()
	return lastInsertId, err
}

func (a *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password, created_at FROM users WHERE email = $1 LIMIT 1`

	var user models.User
	err := a.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		slog.Error("Error getting user", "error", err)
		if err == sql.ErrNoRows {
			return nil, NoRecordsFound
		}
		return nil, err
	}
	return &user, nil
}
