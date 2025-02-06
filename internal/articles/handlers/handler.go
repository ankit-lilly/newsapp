package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/ankibahuguna/newsapp/internal/articles/repository"
	"github.com/ankibahuguna/newsapp/internal/articles/views"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	"github.com/ankibahuguna/newsapp/pkg/errors"
	shared "github.com/ankibahuguna/newsapp/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/ollama/ollama/api"
)

type ArticleService interface {
	GetAllArticles() ([]repository.Article, error)
	GetFeed(category string) ([]repository.Article, error)
	GetOnionFeed() ([]repository.Article, error)
	GetArticleDetail(id int64) (*repository.Article, error)
	GetFavoritesByUser(id int64) ([]repository.Article, error)
	CreateFavoriteArticle(article_id, user_id int64) (*repository.Article, error)
}

func New(a ArticleService, m *melody.Melody, ollamaClient *api.Client) *ArticleHandler {
	return &ArticleHandler{
		ArticleService: a,
		m:              m,
		ollama:         ollamaClient,
	}
}

type ArticleHandler struct {
	ArticleService ArticleService
	m              *melody.Melody
	ollama         *api.Client
}

func (a *ArticleHandler) WebSocketResponse(ctx context.Context, cmp templ.Component, session *melody.Session) error {
	buffer := bytes.Buffer{}
	cmp.Render(ctx, &buffer)
	a.m.BroadcastFilter(buffer.Bytes(), func(q *melody.Session) bool {
		return q.Request.URL.Path == session.Request.URL.Path && q.Request.Header.Get("Sec-WebSocket-Key") == session.Request.Header.Get("Sec-WebSocket-Key")
	})
	return nil
}

func (a *ArticleHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	ctx := context.WithValue(context.Background(), "currentPath", c.Request().URL.Path)
	return cmp.Render(ctx, c.Response().Writer)
}

func (a *ArticleHandler) GetArticlesFromOnion(c echo.Context) error {

	isAuthorized := c.Get("isAuthorized").(bool)
	htmxRequest := c.Get("htmxRequest").(bool)

	articles, err := a.ArticleService.GetOnionFeed()

	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	sl := views.ShowList("| Home", isAuthorized, shared.Categories, views.List(articles))

	if htmxRequest {
		return a.View(c, views.List(articles))
	}

	return a.View(c, sl)
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
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
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
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	htmxRequest := c.Get("htmxRequest").(bool)

	if htmxRequest {
		return a.View(c, views.List(articles))
	}

	sl := views.ShowList("| Favorites", isAuthorized, shared.Categories, views.List(articles))
	return a.View(c, sl)
}

func (a *ArticleHandler) CreateFavArticle(c echo.Context) error {

	var article_id int64
	var err error

	isAuthorized := c.Get("isAuthorized").(bool)

	if !isAuthorized {
		c.Response().Header().Set("Hx-Redirect", "/auth/login")
		return c.Redirect(http.StatusUnauthorized, "/auth/login")
	}

	article_id, err = strconv.ParseInt(c.FormValue("article_id"), 10, 64)

	user := c.Get("user").(*jwt.Token)

	claims := user.Claims.(*auth.JwtClaims)

	var article *repository.Article
	article, err = a.ArticleService.CreateFavoriteArticle(article_id, claims.Id)

	if isAuthorized {
		user := c.Get("user").(*jwt.Token)

		claims := user.Claims.(*auth.JwtClaims)

		c.Logger().Infof("User: %d, Article: %d", claims.Id, article.User)

		if article.User == claims.Id {
			article.IsFavorite = true
		} else {
			article.IsFavorite = false
		}
	}

	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "something went wrong")
	}

	c.Response().Header().Set("HX-Retarget", "#article-detail")
	return a.View(c, views.Detail("", *article))

}

func (a *ArticleHandler) GetArticleDetail(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	isAuthorized := c.Get("isAuthorized").(bool)

	article, err := a.ArticleService.GetArticleDetail(id)

	if isAuthorized {
		user := c.Get("user").(*jwt.Token)

		claims := user.Claims.(*auth.JwtClaims)

		c.Logger().Infof("User: %d, Article: %d", claims.Id, article.User)

		if article.User == claims.Id {
			article.IsFavorite = true
		} else {
			article.IsFavorite = false
		}
	}

	if err != nil {
		c.Logger().Error(err.Error())
		if errors.IsNotFoundError(err) {
			return c.Redirect(http.StatusMovedPermanently, "/404")
		}
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
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

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	article, err := a.ArticleService.GetArticleDetail(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Internal server error")
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		System: "Summarize the input while maintaining the speaking style of a well-educated, Stanford-educated gym bro—confident, energetic, and to the point. Preserve all key details without adding extra information.",
		Model:  "llama3.2",
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
			c.Logger().Error(err.Error())
			return err
		}
		w.Flush()
		return nil
	}

	err = a.ollama.Generate(ctx, req, respFunc)
	if err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "something went wrong")
	}
	return nil
}
