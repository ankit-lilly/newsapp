package hindu

import (
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
	"maps"
	"slices"
)

const ID = "thehindu"

type TheHinduCom struct {
	sources.Source
}

var categories = map[string]string{
	"National":      "news/national",
	"International": "news/international",
	"Business":      "business",
	"Opinion":       "opinion",
	"Sports":        "sport",
	"Entertainment": "entertainment",
	"Science":       "sci-tech/science",
	"Life & Style":  "life-and-style",
}

func NewTheHinduCom() *TheHinduCom {
	return &TheHinduCom{
		sources.Source{
			ID:         ID,
			Name:       "The Hindu",
			BaseURL:    "https://www.thehindu.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				TitleSelector:    "h1.title",
				ContentSelector:  "div.articlebodycontent",
				SubtitleSelector: "h2.sub-title",
				DateSelector:     "div.update-publish-time",
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{
						"div.related-stories-inline",
						".related-topics-list",
						".comments-shares",
						".share-page",
						"script",
						"button",
					},
					RemoveAfterSelectors: []string{
						"button",
					},
				},
			},
		},
	}
}

func (t *TheHinduCom) FeedURL(category string) string {
	var feedURL string
	if category == "" {
		feedURL = fmt.Sprintf("%s/%s", t.BaseURL, "/feeder/default.rss")

	} else {
		feedURL = fmt.Sprintf("%s/%s/feeder/default.rss", t.BaseURL, category)
	}
	return feedURL
}

func (t *TheHinduCom) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}

func (t *TheHinduCom) IsCategoryValid(category string) bool {
	if category == "" {
		return true
	}
	return slices.Contains(slices.Collect(maps.Values(t.Categories)), category)
}
