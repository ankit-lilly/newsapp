package internal

import (
	"github.com/ankibahuguna/newsapp/internal/articles"
	"github.com/ankibahuguna/newsapp/internal/auth"
	"github.com/ankibahuguna/newsapp/pkg/db"
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

}
