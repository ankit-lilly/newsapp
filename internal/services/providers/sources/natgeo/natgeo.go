package natgeo

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
	"net/http"
	"strings"
)

type NatGeo struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.WebFetcher
}

const ID = "natgeo"

var categories = map[string]string{}

func NewNatGeo() *NatGeo {
	return &NatGeo{
		id:         ID,
		name:       "National Geographic",
		baseURL:    "https://www.nationalgeographic.com",
		categories: categories,
		feed: feed.NewWebFetcher(ID, feed.WebSelectors{
			ArticleWrapper: "div.HomepagePromos__row div.ListItem",
			Title:          "a.AnchorLink",
			Description:    "p",
			Link:           "a.AnchorLink",
			PublishedAt:    "time",
		}),
	}
}

func (t *NatGeo) HasCategories() bool {
	return false
}

func (t *NatGeo) ID() string {
	return t.id
}

func (t *NatGeo) Name() string {
	return t.name
}

func (t *NatGeo) FeedURL(category string) string {
	return t.baseURL
}

func (t *NatGeo) Categories() map[string]string {
	return t.categories
}

func (t *NatGeo) Fetch(category string) ([]models.Article, error) {
	return t.feed.Fetch(t.FeedURL(category))
}

func (t *NatGeo) IsCategoryValid(category string) bool {
	return (category == "")
}

func (t *NatGeo) Parse(url string) (models.Article, error) {

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

	doc.Find("div.PrismArticleBody").Children().Each(func(i int, sel *goquery.Selection) {
		finalHTML.WriteString(formatter.FormatNode(sel))
	})

	title := strings.TrimSpace(doc.Find("div.PrismLeadContainer h1").Text())

  if title == "" {
    title = strings.TrimSpace(doc.Find("title").Text())
  }

	subtleTitle := strings.TrimSpace(doc.Find("div.[data-testid='ContentHeaderHed']").Text())
	publishedAt := strings.TrimSpace(doc.Find("div.dateTime").Text())

	return models.Article{
		Title:       title,
		Content:     finalHTML.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
