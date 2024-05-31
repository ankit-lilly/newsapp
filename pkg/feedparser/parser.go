package feedparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"net/http"
	"strings"
)

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	PublishedAt string `json:"publishedAt"`
}

type FeedFetcher struct {
	parser *gofeed.Parser
}

func NewFeedFetcher() *FeedFetcher {
	return &FeedFetcher{parser: gofeed.NewParser()}
}

func (f *FeedFetcher) Fetch(feedURL string) ([]News, error) {
	feed, err := f.parser.ParseURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var news []News
	for _, item := range feed.Items {
		news = append(news, News{
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(item.Description),
			Link:        strings.TrimSpace(item.Link),
			PublishedAt: strings.TrimSpace(item.Published),
		})
	}
	return news, nil
}

type ArticleParser struct{}

func NewArticleParser() *ArticleParser {
	return &ArticleParser{}
}

func filterOutCommentShareWidget(i int, s *goquery.Selection) bool {
	val, exists := s.Parent().Attr("class")
	return !exists || val != "comments-shares share-page"
}

func (a *ArticleParser) Parse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching article: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing article: %w", err)
	}

	var body strings.Builder
	doc.Find("div.articlebodycontent").Find("p").Not(".related-topics-list").FilterFunction(filterOutCommentShareWidget).Each(func(j int, el *goquery.Selection) {
		body.WriteString(el.Text())
		body.WriteString("\n\n")
	})
	return body.String(), nil
}

func (a *ArticleParser) GetRawArticle(u string) (string, error) {
	fmt.Println("Getting article detail", u)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching article: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error fetching article: received status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing article: %w", err)
	}

	var body strings.Builder
	if strings.Contains(u, "onion") {
		fmt.Println("Onion", u)

		doc.Find("*").Each(func(i int, s *goquery.Selection) {
			if _, exists := s.Attr("style"); exists {
				s.RemoveAttr("style")
				s.RemoveAttr("width")
				s.RemoveAttr("height")
				s.RemoveAttr("class")
			}
		})

		selection := doc.Find("main div").First()
		selection.Find(".share-tools-buttons-top, .js_related-stories-inset, .video-html5, .instream-native-video, .js_lightbox, .lightbox, .js_related-stories-inset-mobile, source").Remove()
		html, err := selection.Html()
		if err != nil {
			return "", fmt.Errorf("failed to get HTML: %w", err)
		}
		return html, nil
	}

	doc.Find("div.articlebodycontent").Find("p").Not(".related-topics-list").FilterFunction(filterOutCommentShareWidget).Each(func(j int, el *goquery.Selection) {
		body.WriteString(fmt.Sprintf("<p class='text-lg mt-4'>%s</p>", el.Text()))
	})
	return body.String(), nil
}
