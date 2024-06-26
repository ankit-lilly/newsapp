package articles

import (
	"database/sql"

	"github.com/ankibahuguna/newsapp/internal/articles/handlers"
	"github.com/ankibahuguna/newsapp/internal/articles/repository"
	"github.com/ankibahuguna/newsapp/internal/articles/services"
	"github.com/labstack/echo/v4"
)

func Routes(a *echo.Group, DB *sql.DB) {
	articleRepository := repository.NewRepository(DB)
	articleService := services.NewArticleService(articleRepository)
	articleHandler := handlers.New(articleService)

	a.GET("", articleHandler.GetArticles)
	a.GET("funny", articleHandler.GetArticlesFromOnion)
	a.GET("favorites", articleHandler.GetFavoriteArticles)
	a.POST("favorites", articleHandler.CreateFavArticle)

	a.GET("category/:category", articleHandler.GetArticles)
	a.GET("articles/detail/:id", articleHandler.GetArticleDetail)
	a.GET("articles/detail/:id/summarise", articleHandler.SummariseArticle)
}
