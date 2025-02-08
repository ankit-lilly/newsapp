package feedparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

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
		doc.Find("*").Each(func(i int, s *goquery.Selection) {
			if _, exists := s.Attr("style"); exists {
				s.RemoveAttr("style")
				s.RemoveAttr("width")
				s.RemoveAttr("height")
				s.RemoveAttr("class")
			}
		})

		selection := doc.Find("article").First()
		selection.Find("aside,.single-post-next-story,.wp-block-mbm-post-terms-conditional-container").Remove()
		html, err := selection.Html()
		if err != nil {
			return "", fmt.Errorf("failed to get HTML: %w", err)
		}
		return "<article class='prose mt-4'>" + html + "</article>", nil
	}

	doc.Find("div.articlebodycontent").Find("p").Not(".related-topics-list").FilterFunction(filterOutCommentShareWidget).Each(func(j int, el *goquery.Selection) {
		body.WriteString(fmt.Sprintf("<p class='text-lg mt-4'>%s</p>", el.Text()))
	})
	return body.String(), nil
}
