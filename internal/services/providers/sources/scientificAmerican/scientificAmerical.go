package scientificamerican

import (
	"fmt"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type ScientificAmerican struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "scientificamerican"

func NewScientificAmerican() *ScientificAmerican {
	return &ScientificAmerican{
		sources.Source{
			ID:         ID,
			Name:       "SciAm",
			BaseURL:    "https://www.scientificamerican.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.article__body-ivA3W",
				TitleSelector:    "title",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "p.article_pub_date-EsKM-",
			},
		},
	}
}

func (t *ScientificAmerican) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "platform/syndication/rss")
	return feedURL
}

func (t *ScientificAmerican) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
