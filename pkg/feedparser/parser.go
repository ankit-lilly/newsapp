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
		title := strings.TrimSpace(item.Title)
		description := strings.TrimSpace(item.Description)
		link := strings.TrimSpace(item.Link)
		publishedAt := strings.TrimSpace(item.Published)
		newsItem := News{Title: title, Description: description, Link: link, PublishedAt: publishedAt}
		news = append(news, newsItem)
	}

	return news, nil
}

type ArticleParser struct {
}

func NewArticleParser() *ArticleParser {
	return &ArticleParser{}
}

// Filters out the comment section part of the page
func filterOutCommentShareWidget(i int, s *goquery.Selection) bool {
	val, exists := s.Parent().Attr("class")

	if exists == true && val == "comments-shares share-page" {
		return false
	}

	return true
}

func (a *ArticleParser) Parse(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error fetching article: %s", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", err
	}
	var body strings.Builder
	doc.Find("div.articlebodycontent").Find("p").Not(".related-topics-list").FilterFunction(filterOutCommentShareWidget).Each(func(j int, el *goquery.Selection) {
		body.WriteString(el.Text())
		body.WriteString("\n\n")
	})
	return body.String(), nil

}

func (a *ArticleParser) GetRawArticle(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error fetching article: %s", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", err
	}
	var body strings.Builder
	doc.Find("div.articlebodycontent").Find("p").Not(".related-topics-list").FilterFunction(filterOutCommentShareWidget).Each(func(j int, el *goquery.Selection) {
		body.WriteString(fmt.Sprintf("<p class='text-gray-300 text-lg mt-4'>%s</p>", el.Text()))
	})
	return body.String(), nil

}
