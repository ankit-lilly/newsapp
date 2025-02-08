package main

import (
	"context"
	"embed"
	"github.com/ankibahuguna/newsapp/cmd"
	"github.com/ankibahuguna/newsapp/pkg/db"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed assets/dist
var staticFiles embed.FS

const (
	SECRET_KEY string = "secret"
	DB_NAME    string = "app_data.db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3500"
	}

	if err := db.Init(DB_NAME); err != nil {
		log.Fatal(err)
	}

	server := cmd.NewServer(staticFiles)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(port); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server:", err)
	}
}
