package articles

import (
	"database/sql"
	"log"

	"github.com/ankit-lilly/newsapp/internal/articles/handlers"
	"github.com/ankit-lilly/newsapp/internal/articles/repository"
	"github.com/ankit-lilly/newsapp/internal/articles/services"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"
)

func Routes(a *echo.Group, DB *sql.DB) {
	m := melody.New()

	ollamaClient, err := api.ClientFromEnvironment()

	if err != nil {
		log.Fatalf("failed to create ollama client: %s", err)
	}

	articleRepository := repository.NewRepository(DB)
	articleService := services.NewArticleService(articleRepository)
	articleHandler := handlers.New(articleService, m, ollamaClient)

	a.GET("", articleHandler.GetArticles)
	a.GET("funny", articleHandler.GetArticlesFromOnion)
	a.GET("favorites", articleHandler.GetFavoriteArticles)
	a.POST("favorites", articleHandler.CreateFavArticle)

	a.GET("category/:category", articleHandler.GetArticles)
	a.GET("articles/detail/:id", articleHandler.GetArticleDetail)
	a.GET("articles/detail/:id/chat", articleHandler.Chat)

	m.HandleMessage(articleHandler.HandleChatMessage)

	a.GET("articles/detail/:id/summarise", articleHandler.SummariseArticle)
}
