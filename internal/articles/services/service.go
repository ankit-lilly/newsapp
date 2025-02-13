package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ankit-lilly/newsapp/internal/articles/repository"
	"github.com/ankit-lilly/newsapp/pkg/feedparser"
)

const (
	defaultCategory = "feeder/default.rss"
	baseURL         = "https://www.thehindu.com"
)

type ArticleService struct {
	ArticleRepository *repository.ArticleRepository
}

func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		ArticleRepository: articleRepo,
	}
}

func (a *ArticleService) GetAllArticles() ([]repository.Article, error) {
	return a.ArticleRepository.GetAllArticles()
}

func (a *ArticleService) GetOnionFeed() ([]repository.Article, error) {

	const feedURL = "https://www.theonion.com/rss"
	feed, err := feedparser.NewFeedFetcher().Fetch(feedURL)
	if err != nil {
		return nil, err
	}

	var articles []repository.Article
	for _, item := range feed {
		articles = append(articles, repository.Article{
			ID:          generateIDFromURL(item.Link),
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Body:        "",
			CreatedAt:   getTimeFromDateTimeString(item.PublishedAt),
		})
	}

	err = a.ArticleRepository.BulkInsertArticles(articles)

	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *ArticleService) GetFeed(category string) ([]repository.Article, error) {
	feedURL := fmt.Sprintf("%v/%v", baseURL, category)
	feed, err := feedparser.NewFeedFetcher().Fetch(feedURL)
	if err != nil {
		return nil, err
	}

	var articles []repository.Article
	for _, item := range feed {
		articles = append(articles, repository.Article{
			ID:          generateIDFromURL(item.Link),
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Body:        "",
			CreatedAt:   getTimeFromDateTimeString(item.PublishedAt),
		})
	}

	err = a.ArticleRepository.BulkInsertArticles(articles)

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *ArticleService) GetArticleDetail(id int64) (*repository.Article, error) {
	art, err := a.ArticleRepository.GetArticleByID(id)

	if err != nil {
		return nil, err
	}

	if art.Body != "" {
		return art, nil
	}

	detail, err := feedparser.NewArticleParser().GetRawArticle(art.Link)

	if err != nil {
		return nil, err
	}

	art.Body = detail

	a.ArticleRepository.UpdateArticleByID(*art)

	return art, nil

}

func (a *ArticleService) GetFavoritesByUser(userid int64) ([]repository.Article, error) {
	return a.ArticleRepository.GetFavoritesByUser(userid)
}

func (a *ArticleService) CreateFavoriteArticle(article_id, user_id int64) (*repository.Article, error) {
	isFavorite, err := a.ArticleRepository.IsFavoriteArticle(article_id, user_id)

	if err != nil {

		return nil, err
	}

	if isFavorite {
		fmt.Println("Article already favorited, deleting it")
		err := a.ArticleRepository.DeleteFavoriteArticle(article_id, user_id)

		if err != nil {
			return nil, err
		}

		return a.ArticleRepository.GetArticleByID(article_id)
	}

	fmt.Println("Creating favorite article")
	err = a.ArticleRepository.CreateFavoriteArticle(article_id, user_id)

	if err != nil {
		return nil, err
	}

	return a.ArticleRepository.GetArticleByID(article_id)
}

func generateIDFromURL(url string) int {
	extractID := func(pattern string) (int, error) {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return strconv.Atoi(matches[1])
		}
		return 0, fmt.Errorf("no match found")
	}

	var id int
	var err error

	if strings.Contains(url, "onion") {
		id, err = extractID(`(\d+)$`)
		if err == nil {
			return id
		}
	}

	id, err = extractID(`article(\d+)\.ece`)
	if err == nil {
		return id
	}

	return int(time.Now().Unix())
}

func ConverDateTime(tz string, dt time.Time) string {
	loc, _ := time.LoadLocation(tz)

	return dt.In(loc).Format(time.RFC822Z)
}

func getTimeFromDateTimeString(dateTime string) time.Time {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	parsedTime, err := time.Parse(layout, dateTime)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Now()
	}
	return parsedTime
}
