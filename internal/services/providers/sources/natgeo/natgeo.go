package natgeo

import (
	"strings"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type NatGeo struct {
	sources.Source
}

const ID = "natgeo"

var categories = map[string]string{}

func NewNatGeo() *NatGeo {
	return &NatGeo{
		sources.Source{
			ID:         ID,
			Name:       "Nat Geo",
			BaseURL:    "https://www.nationalgeographic.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: "div.HomepagePromos__row div.ListItem",
				Title:          "a.AnchorLink",
				Description:    "p",
				Link:           "a.AnchorLink",
				PublishedAt:    "time",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.PrismArticleBody",
				TitleSelector:    "div.PrismLeadContainer h1",
				SubtitleSelector: "div.[data-testid='ContentHeaderHed']",
				DateSelector:     "div.dateTime",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"button", ".natgeo-ad"},
				},
			},
		},
	}
}

func (t *NatGeo) FeedURL(category string) string {
	return t.BaseURL
}

func (t *NatGeo) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
