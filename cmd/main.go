package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ankibahuguna/newsapp/internal"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	"github.com/ankibahuguna/newsapp/pkg/db"
	"github.com/ankibahuguna/newsapp/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func cacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/assets/dist/css/") || strings.HasPrefix(c.Request().URL.Path, "/assets/js/") {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		}
		return next(c)
	}
}

func earlyHints(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Link", "</assets/dist/css/style.css>; rel=preload; as=style")
		c.Response().WriteHeader(http.StatusEarlyHints)

		// if f, ok := c.Response().Writer.(http.Flusher); ok {
		// 	f.Flush()
		// }
		return next(c)
	}
}

func IsHTMXRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		htmxHeader := c.Request().Header.Get("HX-Request")
		fmt.Println("Htmx Request here", htmxHeader)

		if htmxHeader == "true" {
			c.Set("htmxRequest", true)
		} else {
			c.Set("htmxRequest", false)
		}
		return next(c)
	}
}

func main() {

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Pre(earlyHints)

	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	e.Use(cacheControl)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtClaims)
		},
		TokenLookup: auth.TokenLookUpKey,
		SigningKey:  []byte(auth.GetJWTSecret()),
		SuccessHandler: func(c echo.Context) {
			c.Set("isAuthorized", true)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			slog.Error(err.Error())
			c.Set("user", nil)
			c.Set("isAuthorized", false)
			return nil
		},
		ContinueOnIgnoredError: true,
	}
	e.Use(echojwt.WithConfig(config))
	e.Use(IsHTMXRequest)

	e.Static("/assets", "assets")

	err := db.Init(DB_NAME)

	if err != nil {
		log.Fatal(err)
	}

	internal.SetUpModules(e)
	// Start Server
	e.Logger.Fatal(e.Start(":5000"))
}
