package feed

import (
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/mmcdole/gofeed"
	"strings"
)

type FeedFetcher struct {
	parser *gofeed.Parser
	portal string
}

func NewFeedFetcher(portal string) *FeedFetcher {
	return &FeedFetcher{parser: gofeed.NewParser(), portal: portal}
}

func (f *FeedFetcher) Fetch(url string) ([]models.Article, error) {

	feed, err := f.parser.ParseURL(url)

	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var article []models.Article

	for _, item := range feed.Items {
		article = append(article, models.Article{
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(item.Description),
			Link:        strings.TrimSpace(item.Link),
			PublishedAt: strings.TrimSpace(item.Published),
			Portal:      f.portal,
		})
	}
	return article, nil
}
