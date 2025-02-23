package fiercepharma

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
)

type FiercePharma struct {
	sources.Source
}

var categories = map[string]string{}

const ID = "fiercepharma"

func NewFiercePharma() *FiercePharma {
	return &FiercePharma{
		sources.Source{
			ID:         ID,
			Name:       "Fierce Pharma",
			BaseURL:    "https://www.fiercepharma.com",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				ContentSelector:  "div#article-body-row",
				TitleSelector:    "h1.element-title",
				SubtitleSelector: "",
				DateSelector:     "span.date",
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{"div.embedded-entity"},
				},
			},
		},
	}
}

func (t *FiercePharma) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, "rss/xml")
	return feedURL
}

func (t *FiercePharma) Fetch(category string) ([]models.Article, error) {
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
