package feed

import (
	"fmt"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/gocolly/colly/v2"
	"log/slog"
)

type WebFetcher struct {
	portal    string
	selectors WebSelectors
}

type WebSelectors struct {
	ArticleWrapper string
	Title          string
	Description    string
	Link           string
	PublishedAt    string
}

func NewWebFetcher(portal string, selectors WebSelectors) *WebFetcher {
	return &WebFetcher{
		portal:    portal,
		selectors: selectors,
	}
}

func (f *WebFetcher) Fetch(url string) ([]models.Article, error) {
	c := colly.NewCollector()
	var articles []models.Article

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	})

	c.OnHTML(f.selectors.ArticleWrapper, func(e *colly.HTMLElement) {
		article := models.Article{
			Title:       e.ChildText(f.selectors.Title),
			Description: truncate(e.ChildText(f.selectors.Description), 200),
			Link:        e.ChildAttr(f.selectors.Link, "href"),
			PublishedAt: e.ChildText(f.selectors.PublishedAt),
			Portal:      f.portal,
		}
		articles = append(articles, article)
	})

	err := c.Visit(url)
	if err != nil {
		slog.Error("Failed to scrape website", err)
		return nil, fmt.Errorf("failed to scrape website: %w", err)
	}

	return articles, nil
}
