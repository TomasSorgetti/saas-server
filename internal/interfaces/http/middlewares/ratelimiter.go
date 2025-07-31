package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"luthierSaas/internal/infrastructure/cache"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiterConfig struct {
    Rate   limiter.Rate 
    Prefix string    
}

func NewRateLimiterMiddleware(cacheService *cache.Cache, config RateLimiterConfig) (gin.HandlerFunc, error) {
    if cacheService.Client() == nil {
        return nil, fmt.Errorf("redis client is nil")
    }

    store, err := redis.NewStoreWithOptions(cacheService.Client(), limiter.StoreOptions{
        Prefix: config.Prefix,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create redis store: %w", err)
    }

    rateLimiter := limiter.New(store, config.Rate)

    middleware := mgin.NewMiddleware(rateLimiter,
        mgin.WithKeyGetter(func(c *gin.Context) string {
            ip := c.ClientIP()
            log.Printf("Rate limiting for IP: %s, Endpoint: %s", ip, c.FullPath())
            return ip
        }),
        mgin.WithErrorHandler(func(c *gin.Context, err error) {
            log.Printf("Rate limit exceeded for IP: %s, Endpoint: %s, Error: %v", c.ClientIP(), c.FullPath(), err)
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error":   "Límite de tasa excedido",
                "message": "Demasiadas solicitudes, intenta de nuevo más tarde",
            })
            c.Abort()
        }),
    )

    return middleware, nil
}