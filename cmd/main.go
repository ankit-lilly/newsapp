package main

import (
	"log"

	"github.com/ankibahuguna/newsapp/internal"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	"github.com/ankibahuguna/newsapp/pkg/db"
	"github.com/ankibahuguna/newsapp/pkg/errors"
	"github.com/ankibahuguna/newsapp/pkg/middlewares"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func main() {

	e := echo.New()

	e.Use(middleware.Gzip())
	e.Pre(middlewares.EarlyHints)
	e.Use(middlewares.CacheControl)
	e.Use(echojwt.WithConfig(auth.JwtConfig))
	e.Use(middlewares.IsHTMXRequest)

	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler
	e.Static("/assets", "assets")

	err := db.Init(DB_NAME)

	if err != nil {
		log.Fatal(err)
	}

	internal.SetUpModules(e)
	e.Logger.Fatal(e.Start(":5000"))
}
