package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
}

func LoadConfig() (*Config, error) {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using default config")
    }

    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        // Fallback for XAMPP MySQL (no password for root)
        databaseURL = "root:@tcp(localhost:3306)/luthier_sass_db?charset=utf8mb4&parseTime=true&loc=Local"
    }

    return &Config{
        DatabaseURL: databaseURL,
    }, nil
}