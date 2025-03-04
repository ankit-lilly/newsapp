package acmqueue

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
	"github.com/ankit-lilly/newsapp/internal/services/providers/formatter"
	"github.com/ankit-lilly/newsapp/internal/services/providers/sources"
	"strings"
)

type ACMQueue struct {
	sources.Source
}

var categories = map[string]string{
	"All":                     "rss/feeds/queuecontent.xml",
	"Web Development":         "rss/feeds/webdevelopment.xml",
	"Programming Languages":   "rss/feeds/programminglanguages.xml",
	"Distributed Development": "rss/feeds/distributeddevelopment.xml",
	"Distributed Computing":   "rss/feeds/distributedcomputing.xml",
	"Concurrency":             "rss/feeds/concurrency.xml",
	"Databases":               "rss/feeds/databases.xml",
	"Performance":             "rss/feeds/performance.xml",
	"Web Security":            "rss/feeds/websecurity.xml",
	"Bioscience":              "rss/feeds/bioscience.xml",
	"Web Services":            "rss/feeds/webservices.xml",
}

const ID = "acmqueue"

func NewACMQueue() *ACMQueue {
	return &ACMQueue{
		sources.Source{
			ID:         ID,
			Name:       "ACMQ",
			BaseURL:    "https://queue.acm.org",
			Categories: categories,
			Fetcher:    feed.NewRSSFetcher(ID),
			ParseConfig: sources.ParseConfig{
				TitleSelector:    "h1.hidetitle",
				ContentSelector:  ".container",
				SubtitleSelector: "",
				DateSelector:     "div.dateTime",
				ContentFilter: sources.ContentFilter{
					RemoveSelectors: []string{".navbar", "label"},
				},
			},
		},
	}
}

func (t *ACMQueue) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.BaseURL, category)
	return feedURL
}

func (t *ACMQueue) Parse(url string) (models.Article, error) {
	doc, err := t.FetchAndParse(url)

	if err != nil {
		return models.Article{}, err
	}

	title := strings.TrimSpace(doc.Find(t.ParseConfig.TitleSelector).Text())

	if title == "" {
		title = doc.Find("h2").Text()
	}

	if title == "" {
		title = doc.Find("h3").Text()
	}

	var stringbuilder strings.Builder

	doc.Find(".floatLeft").Parent().NextAll().Remove()
	doc.Find(".floatLeft").Parent().Remove()

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			stringbuilder.WriteString(formatter.FormatNode(s))
		}
	})

	return models.Article{
		Title:       title,
		Description: doc.Find(t.ParseConfig.SubtitleSelector).Text(),
		Content:     stringbuilder.String(),
		Link:        url,
	}, nil
}

func (t *ACMQueue) Fetch(category string) ([]models.Article, error) {
	return t.Fetcher.Fetch(t.FeedURL(category))
}
