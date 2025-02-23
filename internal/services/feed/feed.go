package feed

import (
	"github.com/ankit-lilly/newsapp/internal/models"
)

type Feed interface {
	Fetch(url string) ([]models.Article, error)
}
