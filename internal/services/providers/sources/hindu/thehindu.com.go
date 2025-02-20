package hindu

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

const ID = "thehindu"

type TheHinduCom struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{
	"National":      "news/national",
	"International": "news/international",
	"Business":      "business",
	"Opinion":       "opinion",
	"Sports":        "sport",
	"Entertainment": "entertainment",
	"Science":       "sci-tech/science",
	"Life & Style":  "life-and-style",
}

func NewTheHinduCom() *TheHinduCom {
	return &TheHinduCom{
		id:         ID,
		name:       "The Hindu",
		baseURL:    "https://www.thehindu.com",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *TheHinduCom) ID() string {
	return t.id
}

func (t *TheHinduCom) Name() string {
	return t.name
}

func (t *TheHinduCom) FeedURL(category string) string {
	var feedURL string
	if category == "" {
		feedURL = fmt.Sprintf("%s/%s", t.baseURL, "/feeder/default.rss")

	} else {
		feedURL = fmt.Sprintf("%s/%s/feeder/default.rss", t.baseURL, category)
	}
	return feedURL
}

func (t *TheHinduCom) HasCategories() bool {
	return true
}

func (t *TheHinduCom) Categories() map[string]string {
	return t.categories
}

func (t *TheHinduCom) Fetch(category string) ([]models.Article, error) {
	return t.feed.Fetch(t.FeedURL(category))
}

func (t *TheHinduCom) IsCategoryValid(category string) bool {
	if category == "" {
		return true
	}
	return slices.Contains(slices.Collect(maps.Values(t.categories)), category)
}

func (t *TheHinduCom) Parse(url string) (models.Article, error) {

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

	var finalHTML strings.Builder

	doc.Find("button").NextAll().Remove()
	doc.Find("div.related-stories-inline, .related-topics-list,.comments-shares, .share-page, button").Remove()

	doc.Find("div.articlebodycontent").Children().Each(func(j int, el *goquery.Selection) {
		finalHTML.WriteString(formatter.FormatNode(el))
	})

	title := strings.TrimSpace(doc.Find("h1.title").Text())
	subtleTitle := strings.TrimSpace(doc.Find("h2.sub-title").Text())
	publishedAt := strings.TrimSpace(doc.Find("div.update-publish-time").Text())

	return models.Article{
		Title:       title,
		Content:     finalHTML.String(),
		Description: subtleTitle,
		PublishedAt: publishedAt,
	}, nil
}
