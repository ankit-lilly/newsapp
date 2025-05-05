package devto

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
	"golang.org/x/sync/errgroup"
)

type DevTo struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "devto"

func NewDevTo() *DevTo {
	return &DevTo{
		sources.Source{
			ID:         ID,
			Name:       "DevTo",
			BaseURL:    "https://dev.to",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: ".crayons-story",
				Title:          "a.crayons-story__hidden-navigation-link",
				Description:    "div.summary-item__dek",
				Link:           "a.crayons-story__hidden-navigation-link",
				PublishedAt:    "time",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  ".crayons-article__main",
				TitleSelector:    ".crayons-article__header__meta h1",
				SubtitleSelector: "",
				DateSelector:     "div.dateTime",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"svg.highlight-action--fullscreen-on", "svg.highlight-action", "div.ltag__link"},
				},
			},
		},
	}
}

func (t *DevTo) FeedURL(category string) string {
	return fmt.Sprintf("%s/%s", t.BaseURL, category)
}

func (t *DevTo) Fetch(category string) ([]models.Article, error) {

	pages := []string{"top/week", "top/month", "top/year", "top/infinity", "/"}

	var (
		mu      sync.Mutex
		results []models.Article
		eg      errgroup.Group
	)

	for _, category := range pages {
		eg.Go(func() error {
			articles, err := t.Fetcher.Fetch(t.FeedURL(category))

			if err != nil {
				return err
			}

			mu.Lock()
			results = append(results, articles...)
			mu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
