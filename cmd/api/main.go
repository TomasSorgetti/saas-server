package main

import (
	"context"
	"log"
	"luthierSaas/internal/di"
	"luthierSaas/internal/infrastructure/config"
	"luthierSaas/internal/infrastructure/persistance"
	"luthierSaas/internal/interfaces/http/middlewares"
	"luthierSaas/internal/interfaces/http/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    db, err := persistance.NewDatabase(cfg)
    if err != nil {
        log.Fatalf("Error connecting to MySQL: %v", err)
    }
    defer db.Close()

    container, emailService := di.NewContainer(db.DB, cfg)
    go emailService.StartWorker(context.Background())

    r := gin.Default()
    r.SetTrustedProxies([]string{"127.0.0.1"})

    // Configurar CORS
    corsConfig := cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
        ExposeHeaders:    []string{"Set-Cookie"}, // Opcional, para depuración
    }
    r.Use(cors.New(corsConfig))

    r.Use(middlewares.ErrorHandler())
    routes.SetupRoutes(r, db.DB, container)

    log.Println("Servidor escuchando en http://localhost:8080")
    r.Run(":8080")
}