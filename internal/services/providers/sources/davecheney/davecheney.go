package davecheney

//https://dave.cheney.net/feed/atom

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services/feed"
)

type DaveCheney struct {
	id         string
	name       string
	baseURL    string
	categories map[string]string
	feed       *feed.FeedFetcher
}

var categories = map[string]string{}

const ID = "davecheney"

func NewDaveCheney() *DaveCheney {
	return &DaveCheney{
		id:         ID,
		name:       "Dave Cheney Blog",
		baseURL:    "https://dave.cheney.net",
		categories: categories,
		feed:       feed.NewFeedFetcher(ID),
	}
}

func (t *DaveCheney) ID() string {
	return t.id
}

func (t *DaveCheney) HasCategories() bool {
	return false
}

func (t *DaveCheney) Name() string {
	return t.name
}

func (t *DaveCheney) FeedURL(category string) string {
	feedURL := fmt.Sprintf("%s/%s", t.baseURL, "feed/atom")
	return feedURL
}

func (t *DaveCheney) Categories() map[string]string {
	return t.categories
}

func (t *DaveCheney) Fetch(category string) ([]models.Article, error) {
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

func (t *DaveCheney) IsCategoryValid(category string) bool {
	return (category == "")
}

func (t *DaveCheney) Parse(url string) (models.Article, error) {
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

	doc.Find("article > div.entry-content").Children().Each(func(j int, el *goquery.Selection) {
		tag := goquery.NodeName(el)
		var html string

		switch tag {
		case "h2":
			html = fmt.Sprintf("<h2 class='text-2xl mt-4'>%s</h2>", el.Text())
		case "pre":
			innerHTML, _ := el.Html()
			html = fmt.Sprintf("<div class='mockup-code'><pre data-prefix=''><code>%s</code></pre></div>", innerHTML)
		case "blockquote":
			html = fmt.Sprintf("<blockquote class='text-lg mt-4'>%s</blockquote>", el.Text())
		case "img":
			src, _ := el.Attr("src")
			html = fmt.Sprintf("<img src='%s' class='img-fluid mt-4'>", src)
		default:
			innerHTML, _ := el.Html()
			html = innerHTML
		}

		body.WriteString(fmt.Sprintf("<div class='text-lg mt-4'>%s</div>", html))
	})

	title := strings.TrimSpace(doc.Find("h1.entry-title").Text())
	publishedAt := strings.TrimSpace(doc.Find("p.date").Text())

	return models.Article{
		Title:       title,
		Content:     body.String(),
		Description: title,
		PublishedAt: publishedAt,
	}, nil
}
