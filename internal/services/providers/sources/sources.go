package sources

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
	"net/http"
	"strings"
)

type ContentFilter struct {
	RemoveSelectors      []string
	RemoveAfterSelectors []string
	CustomFilter         func(*goquery.Selection)
}

type Source struct {
	ID          string
	Name        string
	BaseURL     string
	Categories  map[string]string
	Fetcher     feed.Fetcher
	ParseConfig ParseConfig
}

type ParseConfig struct {
	ContentSelector  string
	TitleSelector    string
	SubtitleSelector string
	DateSelector     string
	TitleProcessor   func(string) string
	ContentFilter    ContentFilter
}

func (p *Source) applyContentFilters(doc *goquery.Document) {
	filter := p.ParseConfig.ContentFilter

	for _, selector := range filter.RemoveAfterSelectors {
		doc.Find(selector).NextAll().Remove()
	}

	for _, selector := range filter.RemoveSelectors {
		doc.Find(selector).Remove()
	}

	for _, selector := range filter.RemoveAfterSelectors {
		doc.Find(selector).NextAll().Remove()
	}

	// Apply custom filter if provided
	if filter.CustomFilter != nil {
		filter.CustomFilter(doc.Selection)
	}
}

func (p *Source) GetID() string {
	return p.ID
}

func (p *Source) GetName() string {
	return p.Name
}

func (p *Source) GetCategories() map[string]string {
	return p.Categories
}

func (p *Source) HasCategories() bool {
	return len(p.Categories) > 0
}

func (p *Source) IsCategoryValid(category string) bool {

	if p.GetID() == "highscalability" {
		return strings.HasPrefix(category, "page/") || category == ""
	}

	if !p.HasCategories() {
		return category == ""
	}
	for _, v := range p.Categories {
		if v == category {
			return true
		}
	}
	return false
}

func (p *Source) Parse(url string) (models.Article, error) {
	doc, err := fetchAndParse(url)
	if err != nil {
		return models.Article{}, err
	}

	p.applyContentFilters(doc)

	var body strings.Builder
	doc.Find(p.ParseConfig.ContentSelector).First().Children().Each(func(j int, el *goquery.Selection) {
		body.WriteString(formatter.FormatNode(el))
	})

	title := strings.TrimSpace(doc.Find(p.ParseConfig.TitleSelector).Text())

	if p.ParseConfig.TitleProcessor != nil {
		title = p.ParseConfig.TitleProcessor(title)
	}

	verticalLine := strings.Index(title, "|")

	if verticalLine > 0 {
		title = title[:verticalLine]
	}

	return models.Article{
		Title:       title,
		Content:     body.String(),
		Description: strings.TrimSpace(doc.Find(p.ParseConfig.SubtitleSelector).Text()),
		PublishedAt: strings.TrimSpace(doc.Find(p.ParseConfig.DateSelector).Text()),
	}, nil
}

func fetchAndParse(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching article: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching article: received status code %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
