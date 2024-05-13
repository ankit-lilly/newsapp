package main

import (
	"github.com/ankibahuguna/newsapp/internal"
	"github.com/ankibahuguna/newsapp/pkg/db"
	"github.com/ankibahuguna/newsapp/pkg/errorHandler"
	"github.com/labstack/echo/v4"
	"log"
)

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func main() {

	e := echo.New()

	e.HTTPErrorHandler = errorhandler.CustomHTTPErrorHandler

	e.Static("/", "assets")

	err := db.Init(DB_NAME)

	if err != nil {
		log.Fatal(err)
	}

	internal.SetUpModules(e)
	// Start Server
	e.Logger.Fatal(e.Start(":5000"))
}
