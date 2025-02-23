package verge

import (
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type Verge struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "verge"

func NewVerge() *Verge {
	return &Verge{
		sources.Source{
			ID:         ID,
			Name:       "Verge",
			BaseURL:    "https://www.theverge.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.duet--layout--entry-body",
				TitleSelector:    "title",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "div.dateTime",
			},
		},
	}
}

func (t *Verge) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "rss/index.xml")
	return feedURL
}

func (t *Verge) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
