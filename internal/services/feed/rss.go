package feed

import (
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/mmcdole/gofeed"
	"strings"
)

type RSSFetcher struct {
	parser     *gofeed.Parser
	portal     string
	portalName string
}

func NewRSSFetcher(portal string) *RSSFetcher {
	return &RSSFetcher{
		parser: gofeed.NewParser(),
		portal: portal,
	}
}

func (f *RSSFetcher) Fetch(url string) ([]models.Article, error) {
	feed, err := f.parser.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var articles []models.Article
	for _, item := range feed.Items {
		articles = append(articles, models.Article{
			Title:       strings.TrimSpace(item.Title),
			Description: truncate(strings.TrimSpace(item.Description), 200),
			Link:        strings.TrimSpace(item.Link),
			PublishedAt: strings.TrimSpace(item.Published),
			Portal:      f.portal,
		})
	}
	return articles, nil
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
