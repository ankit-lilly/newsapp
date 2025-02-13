package internal

import (
	"context"
	"github.com/ankit-lilly/newsapp/internal/articles"
	"github.com/ankit-lilly/newsapp/internal/auth"
	"github.com/ankit-lilly/newsapp/pkg/db"
	"github.com/ankit-lilly/newsapp/pkg/shared"
	"github.com/ankit-lilly/newsapp/pkg/views/pages"
	"github.com/labstack/echo/v4"
)

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func SetUpModules(e *echo.Echo) {

	ar := e.Group("/")
	au := e.Group("/auth")

	err := db.Init(DB_NAME)

	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	dbInstance := db.GetDB()

	articles.Routes(ar, dbInstance)
	auth.Routes(au, dbInstance)

	e.RouteNotFound("/404", func(c echo.Context) error {
		errorPage := pages.ErrorPage("404 | Not Found", false, shared.Categories, "Not Found")
		ctx := context.WithValue(context.Background(), "currentPath", c.Request().URL.Path)
		return errorPage.Render(ctx, c.Response().Writer)
	})

}
