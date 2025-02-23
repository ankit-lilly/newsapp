package feed

import (
	"github.com/ankit-lilly/newsapp/internal/models"
)

type Fetcher interface {
	Fetch(url string) ([]models.Article, error)
}
