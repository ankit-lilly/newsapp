package handlers

import (
	"context"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BaseHandler struct{}

func (bh *BaseHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	isAuthenticated := c.Get("isAuthorized").(bool)

	ctx := context.WithValue(c.Request().Context(), "isAuthorized", isAuthenticated)
	ctx = context.WithValue(ctx, "currentPath", c.Request().URL.Path)

	c.Echo().Logger.Debugf("isAuthorized: %v", isAuthenticated)
	return cmp.Render(ctx, c.Response().Writer)
}

type PageWrapperComp func(string, templ.Component) templ.Component

type RenderProps struct {
	Title            string
	Component        templ.Component
	WrapperComponent PageWrapperComp
	CacheStrategy    string
}

func (bh *BaseHandler) Render(c echo.Context, props RenderProps) error {

	var htmxRequest string = "false"
	isHtmx, ok := c.Get("htmxRequest").(bool)

	if ok && isHtmx {
		htmxRequest = "true"
	}

	switch props.CacheStrategy {
	case "no-cache":
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	default:
		c.Response().Header().Set("Vary", htmxRequest)
		c.Response().Header().Set("Cache-Control", "private, max-age=60")
	}

	if htmxRequest == "true" {
		return bh.View(c, props.Component)
	}

	return bh.View(c, props.WrapperComponent(props.Title, props.Component))
}

func (bh *BaseHandler) RedirectToLogin(c echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
