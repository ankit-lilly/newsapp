package feed

import (
	"fmt"
  "log/slog"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/gocolly/colly/v2"
)

type WebFetcher struct {
	portal string
	// Define site-specific selectors for different elements
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

  slog.Info("Inside fetch", url)
	c.OnHTML(f.selectors.ArticleWrapper, func(e *colly.HTMLElement) {
    slog.Info("Scraping website", e.ChildText(f.selectors.Title))
		article := models.Article{
			Title:       e.ChildText(f.selectors.Title),
			Description: e.ChildText(f.selectors.Description),
			Link:        e.ChildAttr(f.selectors.Link, "href"),
			PublishedAt: e.ChildText(f.selectors.PublishedAt),
			Portal:      f.portal,
		}
		articles = append(articles, article)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape website: %w", err)
	}

	return articles, nil
}
