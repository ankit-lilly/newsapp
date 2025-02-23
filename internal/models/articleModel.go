package models

import (
	"errors"
	"strings"
	"time"
)

type Article struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Portal      string    `json:"portal"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PublishedAt string    `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      int64     `json:"user_id"`
	IsFavorited bool      `json:"is_favorited"`
}

func (a *Article) Validate() error {
	if strings.TrimSpace(a.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if len(a.Title) > 200 {
		return errors.New("title too long")
	}
	return nil
}

func (a *Article) Synopsis() string {
	if len(a.Content) <= 200 {
		return a.Content
	}
	return a.Content[:200] + "..."
}

func (a *Article) RawContent() string {
	return a.Content
}
