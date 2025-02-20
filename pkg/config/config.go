package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	JwtSecret    string
	IsDev        bool
	AppName      string
	DatabaseURL  string
	APP_PORT     string
	CookieDomain string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error reading .env file, using system environment variables")
	}

	return &Config{
		APP_PORT:     getEnv("APP_PORT", "8080"),
		JwtSecret:    getEnv("JWT_SECRET", "secret"),
		IsDev:        getEnv("IS_DEV", "true") == "true",
		AppName:      getEnv("APP_NAME", "newsapp"),
		DatabaseURL:  getEnv("DATABASE_URL", "articles.db"),
		CookieDomain: getEnv("COOKIE_DOMAIN", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
