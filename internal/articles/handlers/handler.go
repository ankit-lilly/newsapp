package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/ankibahuguna/newsapp/internal/articles/repository"
	"github.com/ankibahuguna/newsapp/internal/articles/views"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	shared "github.com/ankibahuguna/newsapp/pkg/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ArticleService interface {
	GetAllArticles() ([]repository.Article, error)
	GetFeed(category string) ([]repository.Article, error)
	GetArticleDetail(id int) (*repository.Article, error)
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

	if isAuthorized {
		user := c.Get("user").(*jwt.Token)

		claims := user.Claims.(*auth.JwtClaims)

		fmt.Println("user", claims.Id, claims.Name, claims.ExpiresAt)

	}

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
	return a.View(c, sl)
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

	sd := views.ShowDetail("| Home", isAuthorized, shared.Categories, views.Detail(tz, *article))

	return a.View(c, sd)
}
