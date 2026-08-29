package feed

import (
	"github.com/ankit-lilly/newsmaxxing/internal/models"
)

type Fetcher interface {
	Fetch(url string) ([]models.Article, error)
}
