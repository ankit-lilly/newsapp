package signalsAndthreads

import (
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
	"log/slog"
	"strings"
)

type SignalsAndThreads struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "signalsandthreads"

func NewSignalsAndThreads() *SignalsAndThreads {
	return &SignalsAndThreads{
		sources.Source{
			ID:         ID,
			Name:       "Signals & Threads",
			BaseURL:    "https://signalsandthreads.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{

				ArticleWrapper: "li.podcast-preview",
				Title:          "h2",
				Description:    "div.blurb",
				Link:           "a",
				PublishedAt:    "h5",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.content-container",
				TitleSelector:    ".podcast-info.h1",
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

func (t *SignalsAndThreads) FeedURL(category string) string {
	return t.BaseURL
}

func (t *SignalsAndThreads) Fetch(category string) ([]models.Article, error) {
	slog.Info("Fetching articles", "info", category)
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))

	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].Link = t.BaseURL + articles[i].Link
	}

	return articles, nil
}
