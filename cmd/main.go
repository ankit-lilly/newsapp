package main

import (
	"log"
	"strings"

	"github.com/ankibahuguna/newsapp/internal"
	"github.com/ankibahuguna/newsapp/pkg/db"
	"github.com/ankibahuguna/newsapp/pkg/errorHandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func cacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/dist/css/") || strings.HasPrefix(c.Request().URL.Path, "/dist/js/") {
			c.Response().Header().Set("Cache-Control", "public, max-age=31536000")
		}
		return next(c)
	}
}

func main() {

	e := echo.New()

	e.Use(middleware.Gzip())
	e.HTTPErrorHandler = errorhandler.CustomHTTPErrorHandler

	e.Use(cacheControl)

	e.Static("/", "assets")

	err := db.Init(DB_NAME)

	if err != nil {
		log.Fatal(err)
	}

	internal.SetUpModules(e)
	// Start Server
	e.Logger.Fatal(e.Start(":5000"))
}
