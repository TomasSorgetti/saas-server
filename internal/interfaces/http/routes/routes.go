package routes

import (
	"database/sql"
	"log"
	"luthierSaas/internal/di"
	"luthierSaas/internal/interfaces/http/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

func SetupRoutes(r *gin.Engine, db *sql.DB, container *di.Container) {
	generalLimiter, err := middlewares.NewRateLimiterMiddleware(container.CacheService, middlewares.RateLimiterConfig{
        Rate:   limiter.Rate{Period: time.Minute, Limit: 100},
        Prefix: "rate:api:v1:free",
    })
    if err != nil {
        log.Fatalf("Failed to initialize general rate limiter: %v", err)
    }
	
    api := r.Group("/v1", generalLimiter)

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
    // auth routes
    SetupAuthRoutes(api, container.AuthHandler, container.CacheService)

	// user routes
    SetupUserRoutes(api, container.UserHandler)
}