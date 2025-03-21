package softwareengineeringdaily

import (
	"strings"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type SFD struct {
	sources.Source
}

const ID = "sed"

var categories = map[string]string{}

func NewSFD() *SFD {
	return &SFD{
		sources.Source{
			ID:         ID,
			Name:       "SED",
			BaseURL:    "https://softwareengineeringdaily.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: "div.article",
				Title:          "h2",
				Description:    ".article__excerpt",
				Link:           "a",
				PublishedAt:    "time",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "#main .row > :nth-child(2)",
				TitleSelector:    "h1.post__title",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "span.article__meta-author",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"button", ".natgeo-ad"},
				},
			},
		},
	}
}

func (t *SFD) FeedURL(category string) string {
	return t.BaseURL
}

func (t *SFD) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
