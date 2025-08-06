package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DatabaseURL   string
	GoogleOAuth   *oauth2.Config
	AppClientURL  string
}

func LoadConfig() (*Config, error) {
	// Cargar archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default config")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback para XAMPP MySQL
		databaseURL = "root:@tcp(localhost:3306)/luthier_sass_db?charset=utf8mb4&parseTime=true&loc=Local"
	}
	
	appClientURL := os.Getenv("APP_CLIENT_URL")
	if databaseURL == "" {
		// Fallback para XAMPP MySQL
		databaseURL = "root:@tcp(localhost:3306)/luthier_sass_db?charset=utf8mb4&parseTime=true&loc=Local"
	}

	// Configuración de Google OAuth2
	googleOAuth := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Config{
		DatabaseURL: databaseURL,
		GoogleOAuth: googleOAuth,
		AppClientURL: appClientURL,
	}, nil
}