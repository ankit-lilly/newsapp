package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/ankit-lilly/newsapp/internal/models"
	"github.com/ankit-lilly/newsapp/internal/services"
	"github.com/ankit-lilly/newsapp/internal/templates/components/articles"
	"github.com/ankit-lilly/newsapp/internal/templates/components/ui"
	"github.com/ankit-lilly/newsapp/internal/templates/pages"
	"github.com/labstack/echo/v4"
)

type ArticleHandler struct {
	*BaseHandler
	articleService *services.ArticleService
}

func NewArticleHandler(articleService *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

func (h *ArticleHandler) ListByCategory(c echo.Context) error {
	category := c.Param("category")
	portal := c.Param("portal")

	articleList, err := h.articleService.GetAll(c.Request().Context(), category, portal)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	return h.Render(c, RenderProps{
		Title:            "Posts",
		Component:        articles.ArticleList(articleList),
		WrapperComponent: pages.Index,
	})
}

func (h *ArticleHandler) List(c echo.Context) error {
	var (
		portalName  string
		articleList []models.Article
		err         error
	)

	portalName = strings.TrimSpace(c.Param("portal"))

	if portalName == "" {
		articleList, err = h.articleService.GetRandomArticles(c.Request().Context())
	} else {
		articleList, err = h.articleService.GetAll(c.Request().Context(), "", portalName)
	}

	if err != nil {
		return h.Render(c, RenderProps{
			Title:            "Error",
			Component:        ui.ErrorBlock(err.Error()),
			WrapperComponent: pages.Index,
		})
	}

	return h.Render(c, RenderProps{
		Title:            "Posts",
		Component:        articles.ArticleList(articleList),
		WrapperComponent: pages.Index,
	})
}

func (h *ArticleHandler) GetArticleByID(c echo.Context) error {

	encodedLink, portalName, err := h.parseAndValidateIdAndPortal(c)
	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}
	link, err := url.QueryUnescape(encodedLink)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	articleDetail, err := h.articleService.GetArticleById(c.Request().Context(), portalName, link)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	if !h.isAuthorized(c) {
		return h.Render(c, RenderProps{
			Title:            articleDetail.Title,
			Component:        articles.Article(*articleDetail),
			WrapperComponent: pages.Index,
		})
	}

	userId, ok := c.Get("userId").(int64)

	if ok {
		articleId, err := h.articleService.IsFavorite(c.Request().Context(), link, userId)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return h.View(c, ui.ErrorBlock("Error checking if article is favorite"))
		}

		if articleId != 0 {
			articleDetail.IsFavorited = true
			articleDetail.ID = articleId
		} else {
			articleDetail.IsFavorited = false
		}
	}

	return h.Render(c, RenderProps{
		Title:            articleDetail.Title,
		Component:        articles.Article(*articleDetail),
		WrapperComponent: pages.Index,
	})
}

func (h *ArticleHandler) GetArticleSummary(c echo.Context) error {

	encodedLink, portalName, err := h.parseAndValidateIdAndPortal(c)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	link, err := url.QueryUnescape(encodedLink)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	contentChan, errChan := h.articleService.GetArticleSummary(c.Request().Context(), portalName, link)

	w := c.Response()
	w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	w.WriteHeader(http.StatusOK)

	for {
		select {
		case content, ok := <-contentChan:
			if !ok {
				contentChan = nil
				return nil
			}

			if _, err := fmt.Fprintf(w, content); err != nil {
				c.Echo().Logger.Error(err)
				return h.View(c, ui.ErrorBlock(err.Error()))
			}
			c.Response().Flush()
		case err, ok := <-errChan:
			if ok && err != nil {
				c.Echo().Logger.Error(err.Error())
				return h.View(c, ui.ErrorBlock(err.Error()))
			} else {
				errChan = nil
			}
		}

		if contentChan == nil && errChan == nil {
			break
		}

	}

	return nil
}

func (h *ArticleHandler) CreateFavoriteArticle(c echo.Context) error {

	encodedLink, portalName, err := h.parseAndValidateIdAndPortal(c)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	link, err := url.QueryUnescape(encodedLink)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	if !h.isAuthorized(c) {
		return h.RedirectToLogin(c)
	}

	userId, ok := c.Get("userId").(int64)

	if !ok {
		c.Echo().Logger.Error("User id not found in context", ok)
		h.RedirectToLogin(c)
	}

	article, err := h.articleService.CreateFavoriteArticle(c.Request().Context(), portalName, link, userId)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	return h.Render(c, RenderProps{
		Title:         "Favorite Articles",
		Component:     ui.Merge([]templ.Component{articles.Article(*article), ui.SuccessBlock("Article added to favorites")}),
		CacheStrategy: "no-cache",
	})

}

func (h *ArticleHandler) DeleteFavoriteArticle(c echo.Context) error {

	articleId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return h.View(c, ui.ErrorBlock("id is required"))
	}

	if !h.isAuthorized(c) {
		return h.RedirectToLogin(c)
	}

	userId, ok := c.Get("userId").(int64)

	if !ok {
		c.Echo().Logger.Error("User id not found in context", ok)
		h.RedirectToLogin(c)
	}

	article, err := h.articleService.GetFavoriteArticle(c.Request().Context(), articleId, userId)

	if err == nil && article == nil {
		c.Echo().Logger.Error("Article not found", articleId, err, userId, articleId)
		return h.View(c, ui.ErrorBlock("Article not found"))
	}

	if err != nil {
		return h.View(c, ui.ErrorBlock("Error fetching article"))
	}

	err = h.articleService.DeleteFavoriteArticle(c.Request().Context(), articleId, userId)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	return h.Render(c, RenderProps{
		Title:         "Favorite Articles",
		Component:     ui.Merge([]templ.Component{articles.Article(*article), ui.SuccessBlock("Article removed from favorites")}),
		CacheStrategy: "no-cache",
	})
}

func (h *ArticleHandler) ListFavoriteArticles(c echo.Context) error {

	if !h.isAuthorized(c) {
		return h.RedirectToLogin(c)
	}

	userId, ok := c.Get("userId").(int64)

	if !ok {
		c.Echo().Logger.Error("User id not found in context", ok)
		return h.View(c, ui.ErrorBlock("Unauthorized"))
	}

	articleList, err := h.articleService.GetFavoriteArticleByUser(c.Request().Context(), userId)

	if err != nil {
		return h.View(c, ui.ErrorBlock(err.Error()))
	}

	return h.Render(c, RenderProps{
		Title:            "Favorite Articles",
		Component:        articles.ArticleList(articleList),
		WrapperComponent: pages.Index,
		CacheStrategy:    "no-cache",
	})
}

func (h *ArticleHandler) parseAndValidateIdAndPortal(c echo.Context) (string, string, error) {
	encodedLink := strings.TrimSpace(c.Param("id"))
	portalName := strings.TrimSpace(c.Param("portal"))

	if encodedLink == "" {
		return "", "", errors.New("id is required")
	}
	if portalName == "" {
		return "", "", errors.New("portal is required")
	}
	return encodedLink, portalName, nil
}

func (h *ArticleHandler) isAuthorized(c echo.Context) bool {
	isAuthorized, ok := c.Get("isAuthorized").(bool)

	if !ok || !isAuthorized {
		return false
	}

	return true
}
