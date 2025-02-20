package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func CacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/assets/dist/css/") || strings.HasPrefix(c.Request().URL.Path, "/assets/js/") {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		}
		return next(c)
	}
}

func EarlyHints(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Link", "</assets/dist/css/style.css>; rel=preload; as=style")
		c.Response().Header().Add("Link", "</assets/js/main.js>; rel=preload; as=script")
		c.Response().WriteHeader(http.StatusEarlyHints)

		return next(c)
	}
}

func IsHTMXRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		htmxHeader := c.Request().Header.Get("HX-Request")

		if htmxHeader == "true" {
			c.Set("htmxRequest", true)
		} else {
			c.Set("htmxRequest", false)
		}
		c.Set("currentPath", c.Request().URL.Path)
		return next(c)
	}
}
