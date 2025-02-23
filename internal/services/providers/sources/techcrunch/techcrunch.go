package techcrunch

import (
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type Techcrunch struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "techcrunch"

func NewTechcrunch() *Techcrunch {
	return &Techcrunch{
		sources.Source{
			ID:         ID,
			Name:       "Techcrunch",
			BaseURL:    "https://techcrunch.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.wp-block-post-content",
				TitleSelector:    "title",
				SubtitleSelector: "h1.article-hero__title",
				DateSelector:     "div.wp-block-post-date",
			},
		},
	}
}

func (t *Techcrunch) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "feed/")
	return feedURL
}

func (t *Techcrunch) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
