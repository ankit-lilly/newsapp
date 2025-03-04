package hindustan

import (
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

const ID = "hindustantimes"

var categories = map[string]string{
	"Opinion":       "opinion",
	"Technology":    "technology",
	"Entertainment": "entertainment",
}

type HindustanTimes struct {
	sources.Source
}

func NewHindustanTimes() *HindustanTimes {
	return &HindustanTimes{
		sources.Source{
			ID:         ID,
			Name:       "HT",
			BaseURL:    "https://hindustantimes.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.detail",
				TitleSelector:    "h1.hdg1",
				SubtitleSelector: "h2.sortDec",
				DateSelector:     "div.dateTime",
			},
		},
	}
}

func (t *HindustanTimes) FeedURL(category string) string {
	return fmt.Sprintf("%s/feeds/rss/%s/rssfeed.xml", t.BaseURL, category)
}

func (t *HindustanTimes) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
