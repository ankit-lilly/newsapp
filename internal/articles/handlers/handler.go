package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/ankibahuguna/newsapp/internal/articles/repository"
	"github.com/ankibahuguna/newsapp/internal/articles/views"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	shared "github.com/ankibahuguna/newsapp/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ollama/ollama/api"
)

type ArticleService interface {
	GetAllArticles() ([]repository.Article, error)
	GetFeed(category string) ([]repository.Article, error)
	GetArticleDetail(id int) (*repository.Article, error)
	GetFavoritesByUser(id int64) ([]repository.Article, error)
	CreateFavoriteArticle(article_id, user_id int64) error
}

func New(a ArticleService) *ArticleHandler {
	return &ArticleHandler{
		ArticleService: a,
	}
}

type ArticleHandler struct {
	ArticleService ArticleService
}

func (a *ArticleHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (a *ArticleHandler) GetArticles(c echo.Context) error {

	category := c.Param("category")

	isAuthorized := c.Get("isAuthorized").(bool)
	htmxRequest := c.Get("htmxRequest").(bool)

	var defaultCategory = "feeder/default.rss"

	if category == "" {
		category = defaultCategory
	} else {
		category = strings.TrimSpace(category) + "/" + defaultCategory
	}

	articles, err := a.ArticleService.GetFeed(category)

	if err != nil {
		return err
	}

	sl := views.ShowList("| Home", isAuthorized, shared.Categories, views.List(articles))

	if htmxRequest {
		return a.View(c, views.List(articles))
	}

	return a.View(c, sl)
}

func (a *ArticleHandler) GetFavoriteArticles(c echo.Context) error {

	isAuthorized := c.Get("isAuthorized").(bool)

	if !isAuthorized {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "You are not authorized to access this page")
	}

	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(*auth.JwtClaims)

	articles, err := a.ArticleService.GetFavoritesByUser(claims.Id)

	if err != nil {
		fmt.Println("error", err)
		return err
	}

	htmxRequest := c.Get("htmxRequest").(bool)
	if htmxRequest {
		return a.View(c, views.List(articles))
	}
	sl := views.ShowList("| Home", isAuthorized, shared.Categories, views.List(articles))
	return a.View(c, sl)
}

func (a *ArticleHandler) CreateFavArticle(c echo.Context) error {

	var article_id int64
	var err error

	article_id, err = strconv.ParseInt(c.FormValue("article_id"), 10, 64)

	isAuthorized := c.Get("isAuthorized").(bool)

	if !isAuthorized {
		return echo.NewHTTPError(echo.ErrUnauthorized.Code, "You are not autorized to peform this action")
	}

	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(*auth.JwtClaims)

	err = a.ArticleService.CreateFavoriteArticle(article_id, claims.Id)

	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "something went wrong")
	}

	return nil

}

func (a *ArticleHandler) GetArticleDetail(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	isAuthorized := c.Get("isAuthorized").(bool)

	article, err := a.ArticleService.GetArticleDetail(id)

	if err != nil {
		return err
	}

	tz := ""
	if len(c.Request().Header["X-Timezone"]) != 0 {
		tz = c.Request().Header["X-Timezone"][0]
	}

	htmxRequest := c.Get("htmxRequest").(bool)

	if htmxRequest {
		return a.View(c, views.Detail(tz, *article))
	}

	sd := views.ShowDetail("| Home", isAuthorized, shared.Categories, views.Detail(tz, *article))

	return a.View(c, sd)
}

func (a *ArticleHandler) SummariseArticle(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	article, err := a.ArticleService.GetArticleDetail(id)

	if err != nil {
		return err
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "something went wrong")
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		System: "Summarize the input provided concisely, ensuring all important details are included.",
		Model:  "llama3",
		Prompt: fmt.Sprintf("Summarize the following text: %s\n", article.Body),
	}

	ctx := context.Background()

	w := c.Response()
	w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	w.WriteHeader(http.StatusOK)

	respFunc := func(resp api.GenerateResponse) error {
		if resp.Done {
			return nil
		}

		if _, err := fmt.Fprintf(c.Response(), resp.Response); err != nil {
			fmt.Println(err)
			return err
		}
		w.Flush()
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "something went wrong")
	}
	return nil
}
