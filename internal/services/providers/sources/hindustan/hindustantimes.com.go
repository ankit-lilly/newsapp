package hindustan

import (
	"slices"

	"fmt"
	"net/http"
	"strings"

	"maps"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
)

type HindustanTimes struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{
	"Opinion":       "opinion",
	"Technology":    "technology",
	"Entertainment": "entertainment",
}

const ID = "hindustantimes"

func NewHindusTanTimes() *HindustanTimes {
	return &HindustanTimes{
		id:         ID,
		name:       "Hindustan Times",
		baseURL:    "https://hindustantimes.com",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *HindustanTimes) ID() string {
	return t.id
}

func (t *HindustanTimes) HasCategories() bool {
	return true
}

func (t *HindustanTimes) Name() string {
	return t.name
}

func (t *HindustanTimes) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/feeds/rss/%s/rssfeed.xml", t.baseURL, category)
	return feedURL
}

func (t *HindustanTimes) Categories() map[string]string {
	return t.categories
}

func (t *HindustanTimes) Fetch(category string) ([]models.Article, error) {
	return t.feed.Fetch(t.FeedURL(category))
}

func (t *HindustanTimes) IsCategoryValid(category string) bool {
	return slices.Contains(slices.Collect(maps.Values(t.categories)), category)
}

func (t *HindustanTimes) Parse(url string) (models.Article, error) {
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

	var body strings.Builder

	doc.Find("div.detail").Children().Each(func(j int, el *goquery.Selection) {
		body.WriteString(formatter.FormatNode(el))
	})

	title := strings.TrimSpace(doc.Find("h1.hdg1").Text())
	subtleTitle := strings.TrimSpace(doc.Find("h2.sortDec").Text())
	publishedAt := strings.TrimSpace(doc.Find("div.dateTime").Text())

	return models.Article{
		Title:       title,
		Content:     body.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
