package hackernoon

import (
	"strings"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type HackerNoon struct {
	sources.Source
}

const ID = "hackernoon"

var categories = map[string]string{}

func NewHackerNoon() *HackerNoon {
	return &HackerNoon{
		sources.Source{
			ID:         ID,
			Name:       "Hackernoon",
			BaseURL:    "https://hackernoon.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: "article",
				Title:          "h2",
				Description:    "",
				Link:           "a",
				PublishedAt:    "small.date",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.profile + div",
				TitleSelector:    "h1.story-title",
				SubtitleSelector: "",
				DateSelector:     "div.dateTime",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"button", "a[href='#commentSection']"},
				},
			},
		},
	}
}

func (t *HackerNoon) FeedURL(category string) string {
	return t.BaseURL + "/c"
}

func (t *HackerNoon) Fetch(category string) ([]models.Article, error) {
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))

	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].Link = t.BaseURL + articles[i].Link
	}

	return articles, nil
}
