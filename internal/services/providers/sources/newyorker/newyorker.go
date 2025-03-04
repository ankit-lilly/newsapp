package newyorker

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type Newyorker struct {
	sources.Source
}

const ID = "newyorker"

var categories = map[string]string{}

func NewNewyorker() *Newyorker {
	return &Newyorker{
		sources.Source{
			ID:         ID,
			Name:       "Newyorker",
			BaseURL:    "https://www.newyorker.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: ".summary-item__content",
				Title:          "h3",
				Description:    "div.summary-item__dek",
				Link:           "a",
				PublishedAt:    "time",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "[class^='ArticlePageChunksContent']",
				TitleSelector:    "h1",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "div.dateTime",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{".journey-unit__container"},
				},
			},
		},
	}
}

func (t *Newyorker) FeedURL(category string) string {
	randomPage := rand.Intn(8) + 1
	if randomPage == 1 {
		return fmt.Sprintf("%s/books/flash-fiction", t.BaseURL)
	}
	return fmt.Sprintf("%s/books/flash-fiction?page=%d", t.BaseURL, randomPage)
}

func (t *Newyorker) Fetch(category string) ([]models.Article, error) {
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))

	if err != nil {
		return nil, err
	}

	for i, article := range articles {
		if strings.Contains(article.Link, "https://www.newyorker.com/") {
			articles[i].Link = article.Link
		} else {
			articles[i].Link = t.BaseURL + article.Link
		}
	}

	return articles, nil

}
