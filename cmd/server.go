package cmd

import (
	"context"
	"embed"
	"net/http"
	"os"

	"github.com/ankibahuguna/newsapp/internal"
	"github.com/ankibahuguna/newsapp/pkg/auth"
	"github.com/ankibahuguna/newsapp/pkg/errors"
	"github.com/ankibahuguna/newsapp/pkg/middlewares"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ziflex/lecho/v3"
)

type Server struct {
	echo   *echo.Echo
	logger *lecho.Logger
}

func NewServer(staticFiles embed.FS) *Server {
	e := echo.New()
	logger := lecho.New(
		os.Stdout,
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)

	e.Logger = logger
	e.Use(lecho.Middleware(lecho.Config{Logger: logger}))
	e.Use(middleware.Gzip())
	e.Use(middlewares.CacheControl)
	e.Use(echojwt.WithConfig(auth.JwtConfig))
	e.Use(middlewares.IsHTMXRequest)
	e.HTTPErrorHandler = errors.CustomHTTPErrorHandler

	e.GET("/assets/*", echo.WrapHandler(http.FileServer(http.FS(staticFiles))))
	internal.SetUpModules(e)

	return &Server{
		echo:   e,
		logger: logger,
	}
}

func (s *Server) Start(port string) error {
	return s.echo.Start(":" + port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
