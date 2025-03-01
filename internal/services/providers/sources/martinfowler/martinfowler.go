package martinfowler

import (
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type MartinFowler struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "martinfowler"

func NewMartinFowler() *MartinFowler {
	return &MartinFowler{
		sources.Source{
			ID:         ID,
			Name:       "Martin Fowler",
			BaseURL:    "https://martinfowler.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.paperBody",
				TitleSelector:    "main > h1",
				SubtitleSelector: "",
				DateSelector:     "p.date",
			},
		},
	}
}

func (t *MartinFowler) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "feed.atom")
	return feedURL
}

func (t *MartinFowler) Fetch(category string) ([]models.Article, error) {
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))
	if err != nil {
		return nil, err
	}
	return articles, nil
}
