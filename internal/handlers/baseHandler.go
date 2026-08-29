package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/ankit-lilly/newsmaxxing/internal/templates"
	"github.com/ankit-lilly/newsmaxxing/internal/templates/components/ui"
	"github.com/labstack/echo/v5"
)

type BaseHandler struct{}

type PageWrapperComp func(string, templ.Component) templ.Component

type RenderProps struct {
	Title            string
	Component        templ.Component
	WrapperComponent PageWrapperComp
	CacheStrategy    string
	CacheDuration    int64
	StatusCode       int
}

func (bh *BaseHandler) View(c *echo.Context, cmp templ.Component) error {
	return bh.view(c, cmp, 0)
}

func (bh *BaseHandler) view(c *echo.Context, cmp templ.Component, statusCode int) error {

	isAuthenticated := c.Get("isAuthorized").(bool)

	ctx := context.WithValue(c.Request().Context(), "isAuthorized", isAuthenticated)
	ctx = context.WithValue(ctx, "currentPath", c.Request().URL.Path)

	c.Logger().Debug("authorization status", "authorized", isAuthenticated)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Header().Set("Vary", "Accept-Encoding, HX-Request")
	if statusCode != 0 {
		c.Response().WriteHeader(statusCode)
	}
	return cmp.Render(ctx, c.Response())
}

func (bh *BaseHandler) Render(c *echo.Context, props RenderProps) error {
	isHtmx, ok := c.Get("htmxRequest").(bool)
	switch props.CacheStrategy {
	case "no-cache":
		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	default:
		if props.CacheDuration > 0 {
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", props.CacheDuration))
		} else {
			c.Response().Header().Set("Cache-Control", "private, max-age=1500")
		}
	}

	if ok && isHtmx {
		return bh.view(c, props.Component, props.StatusCode)
	}
	return bh.view(c, props.WrapperComponent(props.Title, props.Component), props.StatusCode)
}

func (bh *BaseHandler) RedirectToLogin(c *echo.Context) error {
	c.Response().Header().Set("Location", "/login")
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (bh *BaseHandler) RenderError(c *echo.Context, err error) error {
	return bh.Render(c, RenderProps{
		Title:            "Error",
		Component:        ui.ErrorBlock(err.Error()),
		WrapperComponent: templates.Index,
		CacheStrategy:    "no-cache",
	})
}

func (bh *BaseHandler) isAuthorized(c *echo.Context) bool {
	isAuthorized, ok := c.Get("isAuthorized").(bool)

	if !ok || !isAuthorized {
		return false
	}

	return true
}
