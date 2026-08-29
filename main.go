package main

import (
	"context"
	"embed"
	"log"
	"os/signal"
	"syscall"

	"github.com/ankit-lilly/newsmaxxing/cmd"
	"github.com/ankit-lilly/newsmaxxing/pkg/config"
)

//go:embed static/dist
var staticFiles embed.FS

func main() {
	cfg := config.LoadConfig()

	app := cmd.NewApp(cfg)

	if err := app.Init(staticFiles); err != nil {
		log.Fatalf("initialize application: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx, cfg.APP_PORT); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
