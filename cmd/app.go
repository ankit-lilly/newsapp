package cmd

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ankit-lilly/newsmaxxing/internal/db"
	"github.com/ankit-lilly/newsmaxxing/internal/handlers"
	"github.com/ankit-lilly/newsmaxxing/internal/middlewares"
	"github.com/ankit-lilly/newsmaxxing/internal/routes"
	"github.com/ankit-lilly/newsmaxxing/internal/services/llm"
	"github.com/ankit-lilly/newsmaxxing/internal/services/providers"
	"github.com/ankit-lilly/newsmaxxing/pkg/auth"
	"github.com/ankit-lilly/newsmaxxing/pkg/config"
	"github.com/labstack/echo/v5"
	"github.com/ollama/ollama/api"
)

type App struct {
	echo         *echo.Echo
	db           *sql.DB
	ollamaClient *api.Client
	jwtService   *auth.JwtService
	config       *config.Config
}

func NewApp(cfg *config.Config) *App {
	logLevel := slog.LevelInfo
	if cfg.IsDev {
		logLevel = slog.LevelDebug
	}
	e := echo.NewWithConfig(echo.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})),
	})

	if err := db.Init(cfg.DatabaseURL); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
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
		config:       cfg,
	}
}

func (a *App) Start(ctx context.Context, port string) error {
	startConfig := echo.StartConfig{
		Address:         ":" + port,
		GracefulTimeout: 20 * time.Second,
		HideBanner:      true,
	}
	return startConfig.Start(ctx, a.echo)
}

func (a *App) Init(staticFiles embed.FS) error {
	providers.Init()

	errorHandler := handlers.ErrorHandler{BaseHandler: &handlers.BaseHandler{}}
	a.echo.HTTPErrorHandler = errorHandler.CustomHTTPErrorHandler

	authMiddleware := middlewares.NewAuthMiddleware(a.jwtService)

	a.echo.Use(middlewares.CacheControl)
	a.echo.Use(middlewares.IsHTMXRequest)
	a.echo.Use(authMiddleware.JWT())
	a.echo.GET("/static/*", echo.WrapHandler(http.FileServer(http.FS(staticFiles))))

	slog.Info("Using model: ", "info", a.config.ModelToUse)
	llmHandler := llm.New(a.ollamaClient, a.config.ModelToUse)
	routes.RegisterRoutes(a.echo, a.db, llmHandler)

	return nil
}
