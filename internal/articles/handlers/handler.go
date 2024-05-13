package handlers

import (
	"github.com/a-h/templ"
	"github.com/ankibahuguna/newsapp/internal/articles/repository"
	"github.com/ankibahuguna/newsapp/internal/articles/views"
	"github.com/labstack/echo/v4"
	"strconv"
)

type ArticleService interface {
	GetAllArticles() ([]repository.Article, error)
	GetFeed() ([]repository.Article, error)
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

	articles, err := a.ArticleService.GetFeed()

	if err != nil {
		return err
	}

	si := views.ShowList("| Home", views.List(articles))

	return a.View(c, si)
}

func (a *ArticleHandler) GetArticleDetail(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	article, err := a.ArticleService.GetArticleDetail(id)

	if err != nil {
		return err
	}

	tz := ""
	if len(c.Request().Header["X-Timezone"]) != 0 {
		tz = c.Request().Header["X-Timezone"][0]
	}

	si := views.ShowDetail("| Home", views.Detail(tz, *article))

	return a.View(c, si)
}
