package davecheney

//https://dave.cheney.net/feed/atom

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type DaveCheney struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "davecheney"

func NewDaveCheney() *DaveCheney {
	return &DaveCheney{
		sources.Source{
			ID:         ID,
			Name:       "Dave Cheney Blog",
			BaseURL:    "https://dave.cheney.net",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div.entry-content",
				TitleSelector:    "h1.entry-title",
				SubtitleSelector: "",
				DateSelector:     "time.entry-date",
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"div.embedded-entity"},
				},
			},
		},
	}
}

func (t *DaveCheney) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "feed/atom")
	return feedURL
}

func (t *DaveCheney) Fetch(category string) ([]models.Article, error) {
	articles, err := t.Fetcher.Fetch(t.FeedURL(category))
	if err != nil {
		return nil, err
	}

	//Fiercepharma has title with anchor tags
	for i, article := range articles {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Title))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		articles[i].Title = doc.Text()
	}

	return articles, nil
}
