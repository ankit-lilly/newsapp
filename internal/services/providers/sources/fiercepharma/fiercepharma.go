package fiercepharma

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
)

type FiercePharma struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{}

const ID = "fiercepharma"

func NewFiercePharma() *FiercePharma {
	return &FiercePharma{
		id:         ID,
		name:       "Fierce Pharma",
		baseURL:    "https://www.fiercepharma.com",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *FiercePharma) ID() string {
	return t.id
}

func (t *FiercePharma) HasCategories() bool {
	return false
}

func (t *FiercePharma) Name() string {
	return t.name
}

func (t *FiercePharma) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.baseURL, "rss/xml")
	return feedURL
}

func (t *FiercePharma) Categories() map[string]string {
	return t.categories
}

func (t *FiercePharma) Fetch(category string) ([]models.Article, error) {
	articles, err := t.feed.Fetch(t.FeedURL(category))
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

func (t *FiercePharma) IsCategoryValid(category string) bool {
	return (category == "")
}

func (t *FiercePharma) Parse(url string) (models.Article, error) {
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

	doc.Find(".embedded-entity").Remove()

	doc.Find("div#article-body-row").Children().Each(func(j int, el *goquery.Selection) {
		body.WriteString(formatter.FormatNode(el))
	})

	title := strings.TrimSpace(doc.Find("h1.element-title").Text())
	subtleTitle := strings.TrimSpace(doc.Find("div.[data-testid='ContentHeaderHed']").Text())
	publishedAt := strings.TrimSpace(doc.Find("span.date").Text())

	return models.Article{
		Title:       title,
		Content:     body.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
