package feedparser

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"strings"
)

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	PublishedAt string `json:"publishedAt"`
}

type FeedFetcher struct {
	parser *gofeed.Parser
}

func NewFeedFetcher() *FeedFetcher {
	return &FeedFetcher{parser: gofeed.NewParser()}
}

func (f *FeedFetcher) Fetch(feedURL string) ([]News, error) {
	feed, err := f.parser.ParseURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var news []News
	for _, item := range feed.Items {
		news = append(news, News{
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(item.Description),
			Link:        strings.TrimSpace(item.Link),
			PublishedAt: strings.TrimSpace(item.Published),
		})
	}
	return news, nil
}
