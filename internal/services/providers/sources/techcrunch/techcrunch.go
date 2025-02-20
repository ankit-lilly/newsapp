package techcrunch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
)

type Techcrunch struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{}

const ID = "techcrunch"

func NewTechcrunch() *Techcrunch {
	return &Techcrunch{
		id:         ID,
		name:       "Techcrunch",
		baseURL:    "https://techcrunch.com",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *Techcrunch) HasCategories() bool {
	return false
}

func (t *Techcrunch) ID() string {
	return t.id
}

func (t *Techcrunch) Name() string {
	return t.name
}

func (t *Techcrunch) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.baseURL, "feed/")
	return feedURL
}

func (t *Techcrunch) Categories() map[string]string {
	return t.categories
}

func (t *Techcrunch) Fetch(category string) ([]models.Article, error) {
	return t.feed.Fetch(t.FeedURL(category))
}

func (t *Techcrunch) IsCategoryValid(category string) bool {
	return (category == "")
}

func (t *Techcrunch) Parse(url string) (models.Article, error) {

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

	doc.Find("div.wp-block-post-content").Each(func(j int, el *goquery.Selection) {
		body.WriteString(formatter.FormatNode(el))
	})

	title := strings.TrimSpace(doc.Find("title").Text())
	subtleTitle := strings.TrimSpace(doc.Find("h1.article-hero__title").Text())
	publishedAt := strings.TrimSpace(doc.Find("div.wp-block-post-date").Text())

	verticalLine := strings.Index(title, "|")

	if verticalLine > 0 {
		title = title[:verticalLine]
	}

	return models.Article{
		Title:       title,
		Content:     body.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
