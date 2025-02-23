package services

import (
	"context"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/prompts"
	"github.com/ankit-lilly/newsapp/internal/repositories"
	"github.com/ankit-lilly/newsapp/internal/services/llm"
	"github.com/ankit-lilly/newsapp/internal/services/providers"
	"github.com/ollama/ollama/api"
)

type ArticleService struct {
	articleRepo repositories.ArticleRepository
	provider    providers.Providers
	llm         *llm.LLMHandler
}

func NewArticleService(articleRepo repositories.ArticleRepository, llm *llm.LLMHandler, provider providers.Providers) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		llm:         llm,
		provider:    provider,
	}
}

func (s *ArticleService) GetArticleById(ctx context.Context, portalName, id string) (*models.Article, error) {

	portal, err := providers.Get(portalName)

	if err != nil {
		return nil, err
	}

	article, err := portal.Parse(id)

	if err != nil {
		return nil, err
	}

	return &models.Article{
		Title:   article.Title,
		Content: article.Content,
		Portal:  portalName,
		Link:    id,
	}, nil

}

func (s *ArticleService) GetArticleSummary(ctx context.Context, portalName, id string) (<-chan string, <-chan error) {
	article, err := s.GetArticleById(ctx, portalName, id)

	errChan := make(chan error, 1)

	if err != nil {
		errChan <- err
		return nil, errChan
	}

	prompt := fmt.Sprintf("Summarize the following article: %q", article.Content)
	return s.llm.GenerateRequest(ctx, prompts.SUMMARY, prompt, true)
}

func (s *ArticleService) SendChatRequest(ctx context.Context, history []api.Message) (api.Message, error) {

	errorChan, messageChat := s.llm.ChatRequest(ctx, history)

	var (
		resp api.Message
		err  error
	)

	for {
		select {
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			if err != nil {
				err = fmt.Errorf("Error: %v", err)
			}
		case msg, ok := <-messageChat:
			if !ok {
				messageChat = nil
				continue
			}

			resp = msg
			return resp, err
		}

		if messageChat == nil && errorChan == nil {
			return resp, err
		}
	}

}

func (s *ArticleService) GetAll(ctx context.Context, category, portalName string) ([]models.Article, error) {
	portal, err := providers.Get(portalName)

	if err != nil {
		return nil, err
	}

	if !portal.IsCategoryValid(category) {
		return nil, fmt.Errorf("category %q is not valid for portal %q", category, portalName)
	}

	return portal.Fetch(category)
}

func (s *ArticleService) GetRandomArticles(ctx context.Context) ([]models.Article, error) {

	var (
		mu      sync.Mutex
		results []models.Article
		eg      errgroup.Group
	)

	for _, provider := range providers.Registry {
		if provider.GetID() == "natgeo" {
			continue
		}
		provider := provider
		eg.Go(func() error {
			var category string
			if provider.HasCategories() {
				category = slices.Collect(maps.Values(provider.GetCategories()))[rand.Intn(len(provider.GetCategories()))]
			} else {
				category = ""
			}

			articles, err := s.GetAll(ctx, category, provider.GetID())

			if err != nil {
				return err
			}

			mu.Lock()

			for i := len(articles) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				articles[i], articles[j] = articles[j], articles[i]
			}

			articles = articles[:rand.Intn(len(articles))]
			results = append(results, articles...)
			mu.Unlock()
			return nil
		})

	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for i := len(results) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		results[i], results[j] = results[j], results[i]
	}
	return results, nil
}

func (s *ArticleService) CreateFavoriteArticle(ctx context.Context, portalName, link string, userId int64) (*models.Article, error) {

	article, err := s.GetArticleById(ctx, portalName, link)

	if err != nil {
		return nil, err
	}

	article.UserID = userId
	id, err := s.articleRepo.CreateFavoriteArticle(ctx, article)

	if err != nil {
		return nil, err
	}
	article.ID = id
	article.IsFavorited = true

	return article, nil
}

func (s *ArticleService) GetFavoriteArticleByUser(ctx context.Context, userId int64) ([]models.Article, error) {
	return s.articleRepo.GetFavoriteArticleByUser(ctx, userId)
}

func (s *ArticleService) GetFavoriteArticle(ctx context.Context, articleId, userId int64) (*models.Article, error) {
	return s.articleRepo.GetFavoriteArticle(ctx, articleId, userId)
}

func (s *ArticleService) DeleteFavoriteArticle(ctx context.Context, articleId, userId int64) error {
	return s.articleRepo.DeleteFavoriteArticle(ctx, articleId, userId)
}

func (s *ArticleService) IsFavorite(ctx context.Context, link string, userId int64) (int64, error) {
	return s.articleRepo.IsFavorite(ctx, link, userId)
}
