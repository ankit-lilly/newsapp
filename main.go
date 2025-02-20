package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankit-lilly/newsapp/cmd"
	"github.com/ankit-lilly/newsapp/pkg/config"
)

//go:embed static/dist
var staticFiles embed.FS

func main() {

	cfg := config.LoadConfig()

	app := cmd.NewApp(cfg)

	app.Init(staticFiles)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Start(cfg.APP_PORT); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server:", err)
	}
}
