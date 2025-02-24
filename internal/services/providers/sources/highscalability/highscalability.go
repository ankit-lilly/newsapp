package highscalability

import (
	"strings"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type HighScalability struct {
	sources.Source
}

const ID = "highscalability"

var categories = map[string]string{}

func NewHighScalability() *HighScalability {
	return &HighScalability{
		sources.Source{
			ID:         ID,
			Name:       "High Scalability",
			BaseURL:    "https://highscalability.com",
			Categories: categories,
			Fetcher: feed.NewWebFetcher(ID, feed.WebSelectors{
				ArticleWrapper: "article.gh-card",
				Title:          "h3.gh-card-title",
				Description:    "p.gh-card-excerpt",
				Link:           "a.gh-card-link",
				PublishedAt:    "time.gh-card-date",
			}),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "section.gh-content",
				TitleSelector:    "h1.gh-article-title",
				SubtitleSelector: "",
				DateSelector:     "div.dateTime",
				TitleProcessor:   func(title string) string { return strings.TrimSpace(title) },
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"button", ".natgeo-ad"},
				},
			},
		},
	}
}

func (t *HighScalability) FeedURL(category string) string {
	return t.BaseURL + "/page/2"
}

func (t *HighScalability) Fetch(category string) ([]models.Article, error) {
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))

	if err != nil {
		return nil, err
	}

	for i := range articles {
		articles[i].Link = t.BaseURL + articles[i].Link
	}

	return articles, nil
}
