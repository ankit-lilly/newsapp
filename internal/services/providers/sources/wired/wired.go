package wired

import (
	"slices"

	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
	"maps"
)

type Wired struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{
	"AI":          "feed/tag/ai/latest/rss",
	"Top":         "feed/rss",
	"Science":     "feed/category/science/latest/rss",
	"Backchannel": "feed/category/backchannel/latest/rss",
	"Ideas":       "feed/category/ideas/latest/rss",
	"Security":    "feed/category/security/latest/rss",
	"Guides":      "feed/tag/wired-guide/latest/rss",
}

const ID = "wired"

func NewWired() *Wired {
	return &Wired{
		id:         ID,
		name:       "Wired",
		baseURL:    "https://www.wired.com",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *Wired) ID() string {
	return t.id
}

func (t *Wired) HasCategories() bool {
	return true
}

func (t *Wired) Name() string {
	return t.name
}

func (t *Wired) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.baseURL, category)
	return feedURL
}

func (t *Wired) Categories() map[string]string {
	return t.categories
}

func (t *Wired) Fetch(category string) ([]models.Article, error) {
	return t.feed.Fetch(t.FeedURL(category))
}

func (t *Wired) IsCategoryValid(category string) bool {
	return slices.Contains(slices.Collect(maps.Values(t.categories)), category)
}

func (t *Wired) Parse(url string) (models.Article, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Article{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return models.Article{}, fmt.Errorf("error fetching article: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Article{}, fmt.Errorf("error fetching article: received status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return models.Article{}, fmt.Errorf("error parsing article: %w", err)
	}

	// Remove unwanted nodes
	doc.Find("div.article__body").Find("table, div.container--body-inner").Remove()

	var finalHTML strings.Builder

	doc.Find("div.article__body").Children().Each(func(i int, sel *goquery.Selection) {
		finalHTML.WriteString(formatter.FormatNode(sel))
	})

	title := strings.TrimSpace(doc.Find("title").Text())
	subtleTitle := strings.TrimSpace(doc.Find("div.[data-testid='ContentHeaderHed']").Text())
	publishedAt := strings.TrimSpace(doc.Find("div.dateTime").Text())

	verticalLine := strings.Index(title, "|")

	if verticalLine > 0 {
		title = title[:verticalLine]
	}

	return models.Article{
		Title:       title,
		Content:     finalHTML.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
