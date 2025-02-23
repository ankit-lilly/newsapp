package cmd

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"

	"github.com/ankit-lilly/newsapp/internal/db"
	"github.com/ankit-lilly/newsapp/internal/handlers"
	"github.com/ankit-lilly/newsapp/internal/middlewares"
	"github.com/ankit-lilly/newsapp/internal/routes"
	"github.com/ankit-lilly/newsapp/internal/services/llm"
	"github.com/ankit-lilly/newsapp/internal/services/providers"
	"github.com/ankit-lilly/newsapp/pkg/auth"
	"github.com/ankit-lilly/newsapp/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/ollama/ollama/api"
)

type App struct {
	echo         *echo.Echo
	db           *sql.DB
	ollamaClient *api.Client
	jwtService   *auth.JwtService
}

func NewApp(cfg *config.Config) *App {
	e := echo.New()
	e.HideBanner = true

	if cfg.IsDev {
		e.Debug = true
	}

	if err := db.Init(cfg.DatabaseURL); err != nil {
		e.Logger.Fatalf("failed to initialize database: %v", err)
	}

	databaseConn := db.GetDB()

	ollamaClient, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("failed to initialize Ollama client: %v", err)
	}

	return &App{
		echo:         e,
		db:           databaseConn,
		ollamaClient: ollamaClient,
		jwtService:   auth.NewJwtService(),
	}
}

func (a *App) Start(port string) error {
	return a.echo.Start(":" + port)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.echo.Shutdown(ctx)
}

func (a *App) Init(staticFiles embed.FS) error {
	providers.Init()

	errorHandler := handlers.ErrorHandler{}
	a.echo.HTTPErrorHandler = errorHandler.CustomHTTPErrorHandler

	authMiddleware := middlewares.NewAuthMiddleware(a.jwtService)

	a.echo.Use(middlewares.CacheControl)
	a.echo.Use(middlewares.IsHTMXRequest)
	a.echo.Use(authMiddleware.JWT())
	a.echo.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(staticFiles))))

  llmHandler := llm.New(a.ollamaClient, "llama3.2:latest")
	routes.RegisterRoutes(a.echo, a.db, llmHandler)

	return nil

}
