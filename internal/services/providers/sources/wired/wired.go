package wired

import (
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type Wired struct {
	sources.Source
}

var categories = map[string]string{
	"AI":          "feed/tag/ai/latest/rss",
	"Top":         "feed/rss",
	"Science":     "feed/category/science/latest/rss",
	"Backchannel": "feed/category/backchannel/latest/rss",
	"Ideas":       "feed/category/ideas/latest/rss",
	"Security":    "feed/category/security/latest/rss",
	"Guides":      "feed/tag/wired-guide/latest/rss",
}

const ID = "wired"

func NewWired() *Wired {
	return &Wired{
		sources.Source{
			ID:         ID,
			Name:       "Wired",
			BaseURL:    "https://www.wired.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				TitleSelector:    "title",
				ContentSelector:  "[class*='ArticlePageChunks']",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "div.dateTime",
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"div.article__body table", "div.article__body div.container--body-inner"},
				},
			},
		},
	}
}

func (t *Wired) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, category)
	return feedURL
}

func (t *Wired) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
