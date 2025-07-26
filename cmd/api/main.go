package main

import (
	"context"
	"log"
	"luthierSaas/internal/di"
	"luthierSaas/internal/infrastructure/config"
	"luthierSaas/internal/infrastructure/persistance"
	"luthierSaas/internal/interfaces/http/middlewares"
	"luthierSaas/internal/interfaces/http/routes"

	"github.com/gin-gonic/gin"
)

func main() {

    // config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    // database
    db, err := persistance.NewDatabase(cfg)
    if err != nil {
        log.Fatalf("Error connecting to MySQL: %v", err)
    }
    defer db.Close()
    
    // dependency injection
    container, emailService := di.NewContainer(db.DB)

    // Iniciar el worker para procesar la cola de correos
	go emailService.StartWorker(context.Background())
    
    // router
    r := gin.Default()
    r.SetTrustedProxies([]string{"127.0.0.1"})

    // middlewares
    r.Use(middlewares.ErrorHandler())

    routes.SetupRoutes(r, db.DB, container)

    r.Run(":8080")
}